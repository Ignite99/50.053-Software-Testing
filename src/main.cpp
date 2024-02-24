#include <iostream>
#include "./fuzzer/fuzzer.h"

int main() {
    std::cout << "Software Testing Project" << std::endl;
    std::cout << "Team Members:" << std::endl;
    std::cout << "- Group Member 1" << std::endl;
    std::cout << "- Group Member 2" << std::endl;
    std::cout << "- Group Member 3" << std::endl;
    std::cout << "- Group Member 4" << std::endl;
    std::cout << "- Group Member 5" << std::endl;

    // Proof of concept
    std::string input = "Hello, world!";
    std::string mutated_input = bit_flip(input);

    std::cout << "Original input: " << input << std::endl;
    std::cout << "Mutated input (bit flip): " << mutated_input << std::endl;

    return 0;
}