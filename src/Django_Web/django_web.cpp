#include <iostream>
#include <curl/curl.h>
#include <csignal>
#include <fstream>
#include <string>

using namespace std;

static size_t write_callback(char *ptr, size_t size, size_t nmemb, void *userp) {
  return size * nmemb;
}

void signal_handler(int signal) {
    if (signal == SIGINT) {
        cout << "Caught Ctrl+C, exiting..." << endl;
        exit(0);
    }
}

int Django_Handler(string url, string request_type) {
    cout << "Django Web has been called!" << endl;

    signal(SIGINT, signal_handler);

    string file_path;
    ofstream output_file;
    CURL *curl;
    CURLcode res;

    while (true) {
        cout << "Enter file path (or 'q' to quit): ";
        getline(cin, file_path);

        if (file_path == "q") {
            return 1;
        }

        output_file.open(file_path, ios::out | ios::trunc);

        if (output_file.is_open()) {
            break;
        } else {
            cerr << "Invalid file path. Please try again." << endl;
            return 1;
        }
    }

    if (file_path != "q") {
        curl_global_init(CURL_GLOBAL_ALL);

        curl = curl_easy_init();

        if (!curl) {
            cerr << "curl_easy_init() failed" << endl;
            return 1;
        }

        curl_easy_setopt(curl, CURLOPT_URL, "https://www.example.com/api/data");
        curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, write_callback);
        curl_easy_setopt(curl, CURLOPT_WRITEDATA, &output_file);

        res = curl_easy_perform(curl);

        if (res != CURLE_OK) {
            cerr << "curl_easy_perform() failed: " << curl_easy_strerror(res) << endl;
            return 1;
        }

        curl_easy_cleanup(curl);
        curl_global_cleanup();

        cout << "GET request sent successfully!" << endl;
    }

    output_file.close();
    return 0;
}