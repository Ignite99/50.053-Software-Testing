#include <iostream>
#include "fuzzer.h"

using namespace std;

string bit_flip(const string& input) {
  random_device rd;
  mt19937 gen(rd());
  uniform_int_distribution<int> dist(0, input.size() - 1);
  int pos = dist(gen);
  char flipped_byte = input[pos] ^ 1;
  return input.substr(0, pos) + flipped_byte + input.substr(pos + 1);
}

string byte_flip(const string& input) {
  random_device rd;
  mt19937 gen(rd());
  uniform_int_distribution<int> dist(0, input.size() - 1);
  int pos = dist(gen);
  char random_byte = static_cast<char>(rand() % 256);
  return input.substr(0, pos) + random_byte + input.substr(pos + 1);
}

string insert(const string& input) {
  random_device rd;
  mt19937 gen(rd());
  uniform_int_distribution<int> dist(0, input.size());
  int pos = dist(gen);
  char random_byte = static_cast<char>(rand() % 256);
  return input.substr(0, pos) + random_byte + input.substr(pos);
}

string delete_byte(const string& input) {
  if (input.empty()) return input;
  random_device rd;
  mt19937 gen(rd());
  uniform_int_distribution<int> dist(0, input.size() - 1);
  int pos = dist(gen);
  return input.substr(0, pos) + input.substr(pos + 1);
}