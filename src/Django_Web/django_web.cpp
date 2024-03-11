#include <iostream>
#include <curl/curl.h>
#include <csignal>
#include <fstream>
#include <string>

using namespace std;

struct curl_slist *headers;

static size_t write_callback(char *ptr, size_t size, size_t nmemb, void *userp) {
  return size * nmemb;
}

int Django_Handler(string url, string request_type, string input_file_path) {
    ofstream output_file;
    CURLcode res;
    CURL *curl;

    cout << "Django Web has been called!" << endl;

    curl_global_init(CURL_GLOBAL_ALL);
    curl = curl_easy_init();
    if (!curl) {
        cerr << "curl_easy_init() failed" << endl;
        return 1;
    }
    curl_easy_setopt(curl, CURLOPT_URL, url.c_str());
    curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, write_callback);
    curl_easy_setopt(curl, CURLOPT_WRITEDATA, &output_file);

    if (request_type == "GET") {
        res = curl_easy_perform(curl);
    } else if (request_type == "POST") {
        cout << "POST IS CALLED!" << input_file_path << endl;

        ifstream input_file(input_file_path);
        string post_data((istreambuf_iterator<char>(input_file)), istreambuf_iterator<char>());
        input_file.close();

        cout << "POST DATA: "<< post_data.c_str() << endl;

        curl_easy_setopt(curl, CURLOPT_POSTFIELDS, post_data.c_str());

        // Setting headers
        headers = NULL;
        headers = curl_slist_append(headers, "Content-Type: application/json");
        curl_easy_setopt(curl, CURLOPT_HTTPHEADER, headers);

        // POST request performed
        res = curl_easy_perform(curl);

        // Clean headers
        curl_slist_free_all(headers);
    } else {
        cerr << "Invalid request type: " << request_type << endl;
        curl_easy_cleanup(curl);
        curl_global_cleanup();
        return 1;
    }

    // Check response
    if (res != CURLE_OK) {
        cerr << "curl_easy_perform() failed: " << curl_easy_strerror(res) << endl;
        return 1;
    }

    curl_easy_cleanup(curl);
    curl_global_cleanup();

    cout << request_type << " request sent successfully!" << endl;

    output_file.close();
    return 0;
}