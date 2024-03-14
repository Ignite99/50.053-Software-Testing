#ifndef FUZZER_H
#define FUZZER_H

#include <random>
#include <string>
#include <nlohmann/json.hpp>

using namespace std;

string bit_flip(const string& input);
string byte_flip(const string& input);
string insert(const string& input);
string delete_byte(const string& input);
void randomizerJSON(nlohmann::json& j);

#endif