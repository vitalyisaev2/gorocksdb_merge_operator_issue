#include "rocksdb/merge_operator.h"
#include "rocksdb/slice.h"

// RealMergeOperator has been borrowed from the outdated StringAppendTESTOperator from rocksdb repo:
// https://github.com/facebook/rocksdb/blob/4.0.fb/utilities/merge_operators/string_append/stringappend2.h
class RealMergeOperator : public rocksdb::MergeOperator
{
  public:
    RealMergeOperator() {};
    ~RealMergeOperator() {};

    // FullMerge just concatenates all values into another value
    virtual bool FullMerge(const rocksdb::Slice &key,
                           const rocksdb::Slice *existing_value,
                           const std::deque<std::string> &operands,
                           std::string *new_value,
                           rocksdb::Logger *logger) const override;

    virtual bool PartialMerge(const rocksdb::Slice &key,
                              const rocksdb::Slice &left_operand,
                              const rocksdb::Slice &right_operand,
                              std::string *new_value,
                              rocksdb::Logger *logger) const override;

    virtual const char *Name() const override;
};