#include <iostream>
#include <sstream>
#include <vector>
#include <algorithm>
#include <cctype>
#include <set>
#include <fstream>
#include <curl/curl.h>
#include "./fuzzer/fuzzer.h"

// Valid request types
std::set<std::string> valid_request_types = {"GET", "POST", "PUT", "HEAD", "DELETE"};
std::set<std::string> valid_projects = {"coap", "django", "ble"};

int main(int argc, char* argv[]) {
    string project_type, url, request_type, response, line, input_file_path;

    // Check for no. of arguments
    if (argc < 5) {
        std::cerr << "Error: Invalid arguments provided." << endl;
        std::cerr << "The expected format is:" << endl;
        std::cerr << std::string(argv[0]) << " <project type> <url> <request_type> <input_file_path>" << endl;
        return 1;
    }

    project_type = argv[1];
    url = argv[2];
    request_type = argv[3];

    // Project_type & Request_type convert to upper case, url im not checking
    std::transform(request_type.begin(), request_type.end(), request_type.begin(), ::toupper);
    std::transform(project_type.begin(), project_type.end(), project_type.begin(), ::toupper);

    // If proejct_type & request_type is correct
    if (!valid_projects.count(project_type)) {
        cerr << "Invalid project type. Please enter coap/django/ble." << endl;
        return 1;
    }
    if (!valid_request_types.count(request_type)) {
        cerr << "Invalid request type. Please enter GET/POST/PUT/HEAD/DELETE." << endl;
        return 1;
    }

    input_file_path = argv[4];
    std::ifstream input_file(input_file_path);
    if (!input_file.is_open()) {
        std::cerr << "Error: Could not open file '" << input_file_path << "'." << endl;
        return 1;
    }

    // This is where u guys call your functions 
    while (std::getline(input_file, line)) {
        std::cout << line << std::endl;
    }

    input_file.close();

    return 0;
}