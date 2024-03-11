#include <iostream>
#include <curl/curl.h>
#include <csignal>
#include <fstream>
#include <string>

using namespace std;

static size_t write_callback(char *ptr, size_t size, size_t nmemb, void *userp) {
  return size * nmemb;
}

int Django_Handler(string url, string request_type, string input_file_path) {
    cout << "Django Web has been called!" << endl;

    ofstream output_file;
    CURL *curl;
    CURLcode res;

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

    output_file.close();
    return 0;
}