#include <iostream>
#include <curl/curl.h>
#include <csignal>
#include <fstream>
#include <string>
#include "../consts.h"

using namespace std;


struct curl_slist *headers;

CURL* curl;
FILE* output_file;

static size_t write_callback(void *ptr, size_t size, size_t nmemb, FILE *stream) {
    size_t written = fwrite(ptr, size, nmemb, stream);
    fprintf(stream, "\n");
    return written;
}

int check_response(CURLcode res, long http_code, string request_type) {
    switch(http_code) {
        case 200:
            cout << request_type << " request suceeded!" << endl;
            cout << "HTTP Status: " << http_code << ", " << HTTP_STATUS_MESSAGES.at(http_code)<< endl;
            return 0;
        case 201:
            cout << request_type << " create request suceeded!" << endl;
            cout << "HTTP Status: " << http_code << ", " << HTTP_STATUS_MESSAGES.at(http_code)<< endl;
            return 0;
        case 202:
            cout << request_type << " accept request suceeded!" << endl;
            cout << "HTTP Status: " << http_code << ", " << HTTP_STATUS_MESSAGES.at(http_code)<< endl;
            return 0;
        default:
            cout << request_type << " request failed!" << endl;
            cerr << "HTTP status code: " << http_code << ", " << HTTP_STATUS_MESSAGES.at(http_code) << endl;
            return 1;
    }

    cerr << "Error not caught within responses!" << endl;
    return 1;
}

void clean_requests(CURL* curl, FILE *output_file) {
    curl_easy_cleanup(curl);
    curl_global_cleanup();
    fclose(output_file);
}

void request_sender(FILE* output_file, CURL* curl, string request_type, string input_file_path) {
    long http_code;
    CURLcode res;

    http_code = 0;

    // GET/POST requests set up
    if (request_type == "GET") {
        res = curl_easy_perform(curl);
        curl_easy_getinfo(curl, CURLINFO_RESPONSE_CODE, &http_code);
        check_response(res, http_code, request_type); 

    } else if (request_type == "POST") {
        ifstream input_file(input_file_path);
        string post_data((istreambuf_iterator<char>(input_file)), istreambuf_iterator<char>());
        input_file.close();

        curl_easy_setopt(curl, CURLOPT_POSTFIELDS, post_data.c_str());

        // Setting headers
        headers = NULL;
        headers = curl_slist_append(headers, "Content-Type: application/json");
        curl_easy_setopt(curl, CURLOPT_HTTPHEADER, headers);

        // POST request performed
        res = curl_easy_perform(curl);

        curl_easy_getinfo(curl, CURLINFO_RESPONSE_CODE, &http_code);
        // Parse status code
        check_response(res, http_code, request_type);
        // Clean headers
        curl_slist_free_all(headers);
    } else {
        cerr << "Invalid request type: " << request_type << endl;
    }
}

void initialise_requests(string url) {
    curl = nullptr;
    output_file = nullptr;

    curl_global_init(CURL_GLOBAL_ALL);
    curl = curl_easy_init();
    if (!curl) {
        cerr << "curl_easy_init() failed" << endl;
    }

    // Setting output file. TODO: Modify this according to the mutation and fuzzer
    output_file = fopen("./src/fuzzing_responses/response.txt", "ab");
    if (!output_file) {
        cerr << "Failed to open output file" << endl;
        curl_easy_cleanup(curl);
        curl_global_cleanup();
    }

    curl_easy_setopt(curl, CURLOPT_URL, url.c_str());
    curl_easy_setopt(curl, CURLOPT_WRITEDATA, output_file);
    curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, write_callback);
}

int Django_Test_Driver(int energy, string url, string request_type, string input_file_path) {
    int accumulated_iterations;
    bool testing_incomplete;

    // Check if the testing is complete, if so break out of while loop
    testing_incomplete = true;

    // Iterations of the given test, iterations 1 because energy starts at 1.
    accumulated_iterations = 1;

    // Initialise all curl requests with the url and output file
    initialise_requests(url);

    while (testing_incomplete) {

        // This will generate all the unique outputs needed before next mutation
        while (accumulated_iterations != energy) {
            accumulated_iterations++;
            request_sender(output_file, curl, request_type, input_file_path);
        }

        testing_incomplete = false;
        
    }

    // Dealloc all memory allocated to file, headers and curl
    clean_requests(curl, output_file);
    return 0;
} 