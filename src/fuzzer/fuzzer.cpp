#include <iostream>
#include "fuzzer.h"

std::string bit_flip(const std::string& input) {
  std::random_device rd;
  std::mt19937 gen(rd());
  std::uniform_int_distribution<int> dist(0, input.size() - 1);
  int pos = dist(gen);
  char flipped_byte = input[pos] ^ 1;
  return input.substr(0, pos) + flipped_byte + input.substr(pos + 1);
}

std::string byte_flip(const std::string& input) {
  std::random_device rd;
  std::mt19937 gen(rd());
  std::uniform_int_distribution<int> dist(0, input.size() - 1);
  int pos = dist(gen);
  char random_byte = static_cast<char>(std::rand() % 256);
  return input.substr(0, pos) + random_byte + input.substr(pos + 1);
}

std::string insert(const std::string& input) {
  std::random_device rd;
  std::mt19937 gen(rd());
  std::uniform_int_distribution<int> dist(0, input.size());
  int pos = dist(gen);
  char random_byte = static_cast<char>(std::rand() % 256);
  return input.substr(0, pos) + random_byte + input.substr(pos);
}

std::string delete_byte(const std::string& input) {
  if (input.empty()) return input;
  std::random_device rd;
  std::mt19937 gen(rd());
  std::uniform_int_distribution<int> dist(0, input.size() - 1);
  int pos = dist(gen);
  return input.substr(0, pos) + input.substr(pos + 1);
}