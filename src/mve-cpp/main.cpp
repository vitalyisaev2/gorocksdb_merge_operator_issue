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
int step(rocksdb::Iterator *it, int keysTotal, int keysCount);
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
    options.OptimizeLevelStyleCompaction();
    //   // create the DB if it's not already present
    //   options.create_if_missing = true;
    options.merge_operator.reset(new RealMergeOperator());

    // open DB
    rocksdb::Status s = rocksdb::DB::Open(options, kDBPath, &db);
    if (!s.ok())
    {
        throw std::runtime_error("can not open database");
    }

    // iterate over whole database several times
    performIteration(db);

    delete db;
}

void performIteration(rocksdb::DB *db)
{
    for (int i = 0; i < 10; ++i)
    {
        std::cout << "Iteration started" << std::endl;
        iterate(db);
        std::cout << "Iteration finished" << std::endl;
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
    auto keysTotal = estimateNumKeys(db);

    auto readOptions = rocksdb::ReadOptions();
    readOptions.tailing = true;
    readOptions.fill_cache = false;

    int keysCount = 0;
    int valueLenSum = 0;
    auto it = db->NewIterator(readOptions);
    for (it->SeekToFirst(); it->Valid(); it->Next())
    {
        keysCount++;
        valueLenSum += step(it, keysTotal, keysCount);
    };

    if (!it->status().ok())
    {
        throw std::runtime_error(it->status().ToString());
    }

    std::cout << "Sum of value length: " << valueLenSum << std::endl;
}

int step(rocksdb::Iterator *it, int keysTotal, int keysCount)
{
    auto _k = it->key().ToString();
    auto _v = it->value().ToString();
    if (keysCount % (keysTotal / 10) == 0)
    {
        std::cout << "Progress: " << int(100*double(keysCount)/double(keysTotal)) << std::endl;
    }
    return it->value().size();
}