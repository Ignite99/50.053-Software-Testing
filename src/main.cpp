#include <iostream>
#include "./fuzzer/fuzzer.h"
using namespace std;

int main() {
    cout << "Software Testing Project" << endl;
    cout << "Team Members:" << endl;
    cout << "- Group Member 1" << endl;
    cout << "- Group Member 2" << endl;
    cout << "- Group Member 3" << endl;
    cout << "- Group Member 4" << endl;
    cout << "- Group Member 5" << endl;

    // Proof of concept
    string input = "Hello, world!";
    string mutated_input = bit_flip(input);

    cout << "Original input: " << input << endl;
    cout << "Mutated input (bit flip): " << mutated_input << endl;

    return 0;
}