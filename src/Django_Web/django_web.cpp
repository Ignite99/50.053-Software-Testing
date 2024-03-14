#include <iostream>
#include <curl/curl.h>
#include <csignal>
#include <fstream>
#include <string>
#include "../consts.h"

using namespace std;

struct curl_slist *headers;

static size_t write_callback(void *ptr, size_t size, size_t nmemb, FILE *stream) {
    size_t written = fwrite(ptr, size, nmemb, stream);
    fprintf(stream, "\n");
    return written;
}

void monitorFile(const std::string& filename) {
    std::ifstream file(filename);
    std::string line;

    while (true) {
        while (std::getline(file, line)) {
            std::cout << "Line: " << line << std::endl;
        }
        if (file.eof()) {
            break;
        }
        file.clear();
        std::this_thread::sleep_for(std::chrono::seconds(1));
    }
}

int check_response_error(CURLcode res, long http_code, string request_type) {
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

int Django_Handler(string url, string request_type, string input_file_path) {
    FILE *output_file;
    CURLcode res;
    CURL *curl;
    string response_body, selenium_filename;
    long http_code;
    int error_status;

    cout << "Django Web has been called!" << endl;

    // Setting up url
    curl_global_init(CURL_GLOBAL_ALL);
    curl = curl_easy_init();
    if (!curl) {
        cerr << "curl_easy_init() failed" << endl;
        return 1;
    }
    curl_easy_setopt(curl, CURLOPT_URL, url.c_str());

    // Setting up dump for responses, and subsequent iterations of fuzzing
    output_file = fopen("./src/fuzzing_responses/response.txt", "ab");
    if (!output_file) {
        cerr << "Failed to open output file" << endl;
        curl_easy_cleanup(curl);
        curl_global_cleanup();
        return 1;
    }
    curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, write_callback);
    curl_easy_setopt(curl, CURLOPT_WRITEDATA, output_file);

    // Setting base value for http_code
    http_code = 0;

    // GET/POST requests set up
    if (request_type == "GET") {
        res = curl_easy_perform(curl);
        curl_easy_getinfo(curl, CURLINFO_RESPONSE_CODE, &http_code);

        // Parse status code
        error_status = check_response_error(res, http_code, request_type);
        if (error_status != 0) {
            return 1;
        } 

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
        error_status = check_response_error(res, http_code, request_type);
        if (error_status != 0) {
            return 1;
        }

        // Selenium filename to track
        selenium_filename = "./src/Django_Web/selenium/selenium_output.txt";
        try {
            std::thread monitorThread(monitorFile, selenium_filename);
            monitorThread.join();
        } catch (const std::exception& e) {
            std::cerr << "Error creating thread: " << e.what() << std::endl;
        }

        // Clean headers
        curl_slist_free_all(headers);

    } else {
        cerr << "Invalid request type: " << request_type << endl;
        curl_easy_cleanup(curl);
        curl_global_cleanup();
        fclose(output_file);
        return 1;
    }

    curl_easy_cleanup(curl);
    curl_global_cleanup();
    fclose(output_file);

    return 0;
}