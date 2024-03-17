#include <iostream>
#include <sstream>
#include <vector>
#include <algorithm>
#include <cctype>
#include <set>
#include <csignal>
#include <fstream>
#include <curl/curl.h>
#include "./fuzzer/fuzzer.h"
#include "./CoAP_Protocol/coap_protocol.h"
#include "./Django_Web/django_web.h"
#include "./BLE_Zephyr/ble_zephyr.h"

// Valid request & project types
std::set<std::string> valid_request_types = {"GET", "POST", "PUT", "HEAD", "DELETE"};
std::set<std::string> valid_projects = {"COAP", "DJANGO", "BLE"};

void signal_handler(int signal) {
    if (signal == SIGINT) {
        cout << "Caught Ctrl+C, exiting..." << endl;
        exit(0);
    }
}

// Call ./Software_Testing_Project.exe
int main(int argc, char* argv[]) {
    string project_type, url, request_type, response, line, input_file_path, output_file_path;
    std::ifstream input_file(input_file_path);
    std::ifstream output_file(output_file_path);

    signal(SIGINT, signal_handler);

    if(argc == 6){
        project_type = argv[1];
        url = argv[2];
        request_type = argv[3];

        // Project_type & Request_type convert to upper case, url im not checking
        std::transform(request_type.begin(), request_type.end(), request_type.begin(), ::toupper);
        std::transform(project_type.begin(), project_type.end(), project_type.begin(), ::toupper);

        // Check if project_type & request_type is correct
        if (!valid_projects.count(project_type)) {
            cerr << "Invalid project type. Please enter COAP/DJANGO/BLE." << endl;
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

        output_file_path = argv[5];
        std::ifstream output_file(output_file_path);
        if (!output_file.is_open()) {
            std::cerr << "Error: Could not open file '" << output_file_path << "'." << endl;
            return 1;
        } 

    } else if(argc == 1){
        // Input project_type
        while(true){
            std::cout << "What fuzzing target are you testing? [Options: COAP, DJANGO, BLE]" << endl;
            std::cin >> project_type;
            std::transform(project_type.begin(), project_type.end(), project_type.begin(), ::toupper);
            if (!valid_projects.count(project_type)) {
                cerr << "Invalid project type. Please enter COAP/DJANGO/BLE.\n" << endl;
            } else {
                break; 
            }
        }
        
        // Input URL
        std::cout << "\nWhat URL are you testing on?" << endl;
        std::cin >> url;
        
        // Input request_type
        while(true){
            std::cout << "\nWhat is your request type? [Options: GET, POST, PUT, HEAD, DELETE]" << endl;
            std::cin >> request_type;
            std::transform(request_type.begin(), request_type.end(), request_type.begin(), ::toupper);
            if (!valid_request_types.count(request_type)) {
                cerr << "Invalid request type. Please enter GET/POST/PUT/HEAD/DELETE.\n" << endl;
            } else {
                break;
            }
        }
        
        // Input input_file_path
        while(true){
            std::cout << "\nWhat is your seed input's file path?" << endl;
            std::cin >> input_file_path;
            std::ifstream input_file(input_file_path);
            if (!input_file.is_open()) {
                std::cerr << "Error: Could not open file '" << input_file_path << "'." << endl;
            } else {
                break;
            }
        }

        // Input output_file_path
        while(true){
            std::cout << "\nWhat is your output file path (output file contains statistics and bug report)?" << endl;
            std::cin >> output_file_path;
            std::ifstream output_file(output_file_path);
            if (!output_file.is_open()) {
                std::cerr << "Error: Could not open file '" << output_file_path << "'." << endl;
            } else {
                break;
            }
        }
        
    } else {
        std::cerr << "Error: Invalid arguments provided." << endl;
        std::cerr << "The expected format is:" << endl;
        std::cerr << std::string(argv[0]) << " <project type> <url> <request_type> <input_file_path> <output_file_path>" << endl;
        return 1;
    }

    // This is where you guys call your functions 
    if (project_type == "COAP") {
        std::cout << "\n[COAP] Fuzzer has initiated call to COAP!" << endl;
        CoAP_Handler();

    } else if (project_type == "DJANGO") {
        std::cout << "\n[DJANGO] Fuzzer has initiated call to DJANGO Web Application!" << endl;
        Django_Test_Driver(1, url, request_type, input_file_path);

    } else if (project_type == "BLE") {
        std::cout << "\n[BLE] Fuzzer has initiated call to BLE_Zephyr!" << endl;
        BLE_Zephyr_Handler();

    } else {
        cerr << "Project type mutated. Project type: " << project_type << ". Check code now!" << endl;
        return 1;
    }

    // Close input file
    input_file.close();
    output_file.close();
    return 0;
}