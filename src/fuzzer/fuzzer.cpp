#include <iostream>
#include <random>
#include <nlohmann/json.hpp>
#include "fuzzer.h"

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

string randomize_string(int length) {
    static const char alphanum[] =
        "0123456789"
        "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
        "abcdefghijklmnopqrstuvwxyz";

    std::random_device rd;
    std::mt19937 gen(rd());
    std::uniform_int_distribution<> distr(0, sizeof(alphanum) - 2);

    std::string random_string;
    random_string.reserve(length);

    for (int i = 0; i < length; ++i) {
        random_string += alphanum[distr(gen)];
    }

    return random_string;
}

/* ============================== DJANGO FUZZER ============================== */

// Django JSON fuzzer
void randomizerJSON(nlohmann::json& j) {
  int num_value;
  string string_value;

  std::random_device rd;
  std::mt19937 gen(rd());
  std::uniform_int_distribution<> int_dist(0, 100);
  std::uniform_real_distribution<> float_dist(0.0, 1.0);
  std::uniform_int_distribution<> bool_dist(0, 1);

  for (auto& [key, value] : j.items()) {
    if (value.is_number()) {
      num_value = int_dist(gen);
      j[key] = num_value;
    } else if (value.is_string()) {
      string_value = randomize_string(50);
      j[key] = string_value;
    }
    std::cout << "Value of field " << key << ": " << value << endl;
  }
}


/* ============================== DJANGO FUZZER ============================== */