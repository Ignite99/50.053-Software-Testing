#ifndef FUZZER_H
#define FUZZER_H

#include <random>
#include <string>

using namespace std;

string bit_flip(const string& input);
string byte_flip(const string& input);
string insert(const string& input);
string delete_byte(const string& input);

#endif