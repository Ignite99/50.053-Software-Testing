#ifndef FUZZER_H
#define FUZZER_H

#include <random>
#include <string>

std::string bit_flip(const std::string& input);
std::string byte_flip(const std::string& input);
std::string insert(const std::string& input);
std::string delete_byte(const std::string& input);

#endif