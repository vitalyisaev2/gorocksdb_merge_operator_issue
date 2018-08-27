#include <iostream>
#include <string>
#include <chrono>
#include <thread>
#include <exception>

#include "rocksdb/db.h"
#include "rocksdb/slice.h"
#include "rocksdb/options.h"

#include "real_merge_operator.h"

std::string kDBPath = "./segments";

void performIteration(rocksdb::DB *db);
void iterate(rocksdb::DB *db);
int estimateNumKeys(rocksdb::DB *db);
void step(rocksdb::Iterator *it, int count);
void run();

int main()
{
    try {
        run();
    } catch (const std::exception& ex) {
        std::cout << "Exception: " << ex.what() << std::endl;
        return 1;
    }
    return 0;
}

void run() {
    rocksdb::DB *db;
    rocksdb::Options options;

    //   // Optimize RocksDB. This is the easiest way to get RocksDB to perform well
    //   options.IncreaseParallelism();
    //   options.OptimizeLevelStyleCompaction();
    //   // create the DB if it's not already present
    //   options.create_if_missing = true;
    auto mo = new RealMergeOperator();
    options.merge_operator.reset(mo);

    // open DB
    rocksdb::Status s = rocksdb::DB::Open(options, kDBPath, &db);
    if (!s.ok())
    {
        throw std::runtime_error("can not open database");
    }

    // iterate over whole database several times
    performIteration(db);

    delete db;
    delete mo;
}

void performIteration(rocksdb::DB *db)
{
    for (int i = 0; i < 5; ++i)
    {
        std::cout << "Iteration started" << std::endl;
        iterate(db);
        std::cout << "Iteration finished" << std::endl;
        std::this_thread::sleep_for(std::chrono::seconds(10));
    }
}

int estimateNumKeys(rocksdb::DB *db)
{
    std::string num;
    auto ok = db->GetProperty("rocksdb.estimate-num-keys", &num);
    if (!ok)
    {
        throw std::runtime_error("can not estimate number of keys in database");
    }

    return std::stoi(num);
}

void iterate(rocksdb::DB *db)
{
    auto totalKeys = estimateNumKeys(db);
    std::cout << totalKeys << std::endl;

    auto readOptions = rocksdb::ReadOptions();
    readOptions.tailing = true;
    readOptions.fill_cache = false;

    int counter = 0;
    auto it = db->NewIterator(readOptions);
    for (it->SeekToFirst(); it->Valid(); it->Next())
    {
        counter++;
        step(it, counter);
    };

    if (!it->status().ok())
    {
        throw std::runtime_error(it->status().ToString());
    }
}

void step(rocksdb::Iterator *it, int counter)
{
    if (counter % 1000000 == 0)
    {
        std::cout << "Progress: " << counter << std::endl;
        std::cout << "Example: " << it->key().ToString() << " " << it->value().ToString() << std::endl;
    }
}