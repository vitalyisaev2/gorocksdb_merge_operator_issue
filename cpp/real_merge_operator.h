#include "rocksdb/merge_operator.h"
#include "rocksdb/slice.h"

class RealMergeOperator : public rocksdb::MergeOperator
{
  public:
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