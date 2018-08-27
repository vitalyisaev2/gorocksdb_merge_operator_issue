#include <iostream>
#include <string>
#include <chrono>
#include <thread>
#include <exception>

#include "rocksdb/db.h"
#include "rocksdb/slice.h"
#include "rocksdb/options.h"

using namespace rocksdb;

std::string kDBPath = "./segments";

void performIteration(DB *db);
void iterate(DB *db);
int estimateNumKeys(DB *db);
void step(Iterator* it, int count);

int main() {
    DB* db;
    Options options;
    
    //   // Optimize RocksDB. This is the easiest way to get RocksDB to perform well
    //   options.IncreaseParallelism();
    //   options.OptimizeLevelStyleCompaction();
    //   // create the DB if it's not already present
    //   options.create_if_missing = true;
    
    // open DB
    Status s = DB::Open(options, kDBPath, &db);
    assert(s.ok());
    
    // iterate several times
    performIteration(db);
    
    delete db;
    return 0;
}

void performIteration (DB* db) {
    for (int i = 0; i < 5; ++i) {
        std::cout << "Iteration started" << std::endl;
        iterate(db);
        std::cout << "Iteration finished" << std::endl;
        std::this_thread::sleep_for(std::chrono::seconds(10));
    }
}

int estimateNumKeys(DB *db) {
    std::string num;
    auto ok = db->GetProperty("rocksdb.estimate-num-keys", &num);
    if (!ok) {
        throw std::runtime_error("can not estimate number of keys in database");
    }

    return std::stoi(num);
}

void iterate(DB* db) {
    auto totalKeys = estimateNumKeys(db);
    std::cout << totalKeys << std::endl;

    auto readOptions = ReadOptions();
    readOptions.tailing = true;
    readOptions.fill_cache = false;

    int counter = 0;
    auto it = db->NewIterator(readOptions);
    for (it->SeekToFirst(); it->Valid(); it->Next()) {
        counter++;
        step(it, counter);
    };

    if (!it->status().ok()) {
        throw std::runtime_error(it->status().ToString());
    }
}

void step(Iterator* it, int counter) {
    if (counter % 1000000 == 0) {
        std::cout << "Progress: " << counter << std::endl;
        std::cout << "Example: " << it->key().ToString() << " " << it->value().ToString() << std::endl;
    }
}