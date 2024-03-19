#include <iostream>
#include <random>
#include <nlohmann/json.hpp>
#include "fuzzer.h"

string bit_flip(const string &input)
{
  random_device rd;
  mt19937 gen(rd());
  uniform_int_distribution<int> dist(0, input.size() - 1);
  int pos = dist(gen);
  char flipped_byte = input[pos] ^ 1;
  return input.substr(0, pos) + flipped_byte + input.substr(pos + 1);
}

string byte_flip(const string &input)
{
  random_device rd;
  mt19937 gen(rd());
  uniform_int_distribution<int> dist(0, input.size() - 1);
  int pos = dist(gen);
  char random_byte = static_cast<char>(rand() % 256);
  return input.substr(0, pos) + random_byte + input.substr(pos + 1);
}

string insert(const string &input)
{
  random_device rd;
  mt19937 gen(rd());
  uniform_int_distribution<int> dist(0, input.size());
  int pos = dist(gen);
  char random_byte = static_cast<char>(rand() % 256);
  return input.substr(0, pos) + random_byte + input.substr(pos);
}

string delete_byte(const string &input)
{
  if (input.empty())
    return input;
  random_device rd;
  mt19937 gen(rd());
  uniform_int_distribution<int> dist(0, input.size() - 1);
  int pos = dist(gen);
  return input.substr(0, pos) + input.substr(pos + 1);
}

string randomize_string(int length)
{
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

  for (int i = 0; i < length; ++i)
  {
    random_string += alphanum[distr(gen)];
  }

  return random_string;
}

/* ============================== DJANGO FUZZER ============================== */

// Django JSON fuzzer

/* ============================== DJANGO FUZZER ============================== */