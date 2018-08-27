#include "real_merge_operator.h"

// implementation is adopted from examples widely spread by github

bool RealMergeOperator::FullMerge(const rocksdb::Slice &key,
                                  const rocksdb::Slice *existing_value,
                                  const std::deque<std::string> &operands,
                                  std::string *new_value,
                                  rocksdb::Logger *logger) const
{
    // clear the *new_value for writing.
    assert(new_value);
    new_value->clear();

    // estimate number of bytes
    int numBytes = 0;
    for (auto it = operands.begin(); it != operands.end(); ++it)
    {
        numBytes += it->size();
    }

    // reserve space
    if (existing_value)
    {
        numBytes += existing_value->size();
        new_value->reserve(numBytes);
        new_value->append(existing_value->data(), existing_value->size());
    }
    else
    {
        new_value->reserve(numBytes);
    }

    // concatenate values
    for (auto it = operands.begin(); it != operands.end(); ++it)
    {
        new_value->append(*it);
    }
    return true;
}

bool RealMergeOperator::PartialMerge(const rocksdb::Slice &key,
                                     const rocksdb::Slice &left_operand,
                                     const rocksdb::Slice &right_operand,
                                     std::string *new_value,
                                     rocksdb::Logger *logger) const
{
    return false;
}

const char *RealMergeOperator::Name() const { return "real"; }