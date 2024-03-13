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
    std::string random_string;
    std::random_device rd;
    std::mt19937 gen(rd());
    
    // Just to map out known letters, we need expand this into non-English and symoblic territory
    static const char alphanum[] =
        "0123456789"
        "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
        "abcdefghijklmnopqrstuvwxyz";

    std::uniform_int_distribution<> distr(0, sizeof(alphanum) - 2);

    /*
    If you are wondering why it lags out when we put such a long string in, its this that causes it
    allocating a string length up to integer limit is quite fucked haha
    */
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
  float float_value;
  bool bool_value;

  std::random_device rd;
  std::mt19937 gen(rd());

  /* 
  Take note, we can attempt a long ass string, but idk how long it will lag the fuzzer
  so ideally i want to go for just a 100, there should be better ways to implement this.
  */
  std::uniform_int_distribution<> string_length_dist(0, 100);

  /*
  Your typical randomising numbers for integer, float and boolean
  */ 
  std::uniform_int_distribution<> int_dist(std::numeric_limits<int>::min(), std::numeric_limits<int>::max());
  std::uniform_real_distribution<> float_dist(std::numeric_limits<float>::min(), std::numeric_limits<float>::max());
  std::uniform_int_distribution<> bool_dist(std::numeric_limits<bool>::min(), std::numeric_limits<bool>::max());

  /*
  We can change it to a switch and abstract. Naive method would be this for now.
  */
  for (auto& [key, value] : j.items()) {
    if (value.is_number()) {
      num_value = int_dist(gen);
      j[key] = num_value;
    } else if (value.is_string()) {
      string_value = randomize_string(string_length_dist(gen));
      j[key] = string_value;
    } else if (value.is_number_float()) {
      float_value = float_dist(gen);
      j[key] = float_value;
    } else if (value.is_boolean()) {
      bool_value = bool_dist(gen);
      j[key] = bool_value;
    } else if (value.is_null()) {
      // Honestly idk how to fuzz this
      num_value = int_dist(gen);
      j[key] = num_value;
    } else if (value.is_array()) {
      // How da hell am I gonna fuzz the array parts haha
      num_value = int_dist(gen);
      j[key] = num_value;
    } else {
      std::cout << "value type is not covered" << endl;
    }


    /*
    
    TODO pls add if you guys think of more: 
    1) Add detecting type of data struct and switching to a different data type
    2) Add array catch length randomizer and type detector

    */

    // Template print statement
    std::cout << "{ " << key << ", " << value << " }" << endl;
  }
}


/* ============================== DJANGO FUZZER ============================== */