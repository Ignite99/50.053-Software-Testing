#include <iostream>
#include <curl/curl.h>
#include <csignal>
#include <fstream>
#include <sstream>
#include <string>
#include <list>
#include <tuple>
#include <nlohmann/json.hpp>
#include <ctime>
#include "../consts.h"
#include "../fuzzer/fuzzer.h"
#include "../fuzzer/json_mutator/json_mutator.h"
#include "../HTML_Logger/html_logger.h"

using namespace std;
using json = nlohmann::json;

struct curl_slist *headers;

HTMLLogger html_logger("./src/fuzzing_responses/", "logs.html", "DJANGO");

typedef struct json_seed
{
    json data;
    string key_to_mutate;
    int energy;
} json_seed;

CURL *curl;
FILE *output_file;

json parse_json(string input_file_path)
{
    json json_data;

    ifstream f(input_file_path);
    json data = json::parse(f);

    return data;
}

string get_last_line()
{
    string filename;
    string lastLine;
    ifstream fin;
    char ch;

    // TODO: remove hardcode for the filename
    filename = "./src/fuzzing_responses/response.txt";

    fin.open(filename);
    if (fin.is_open())
    {
        // Start reading from end of file
        fin.seekg(-2, ios_base::end);

        bool keepLooping = true;
        while (keepLooping)
        {
            // Read one char from file
            fin.get(ch);

            if ((int)fin.tellg() <= 1)
            {
                fin.seekg(0);
                keepLooping = false;
            }
            else if (ch == '\n')
            {
                keepLooping = false;
            }
            else
            {
                fin.seekg(-2, ios_base::cur);
            }
        }
        getline(fin, lastLine);

        fin.close();
        return lastLine;
    }
    return 0;
}

bool is_interesting(string &line, long http_code)
{
    if (line.find("\"success\": false"))
    {
        return true;
    }
    if (http_code != 200)
    {
        return true;
    }

    // TODO: maybe can check if previous rows are of the same type
    // - but this proves to be hard cause its not really flexible

    return false;
}

static size_t write_callback(void *ptr, size_t size, size_t nmemb, FILE *stream)
{
    size_t written = fwrite(ptr, size, nmemb, stream);
    fprintf(stream, "\n");
    return written;
}

int check_response(CURLcode res, long http_code, string request_type, string body)
{
    vector<string> row;
    time_t now;
    char *dt;

    // current date/time based on current system
    now = time(0);

    // convert now to string form
    dt = ctime(&now);

    row = {dt, request_type, body, to_string(http_code)};

    // log responses in html_logger after printing out status message
    switch (http_code)
    {
    case 200:
        cout << request_type << " request suceeded!" << endl;
        cout << "HTTP Status: " << http_code << ", " << HTTP_STATUS_MESSAGES.at(http_code) << endl;
        html_logger.add_row("background-color:palegreen", row);
        return 0;
    case 201:
        cout << request_type << " create request suceeded!" << endl;
        cout << "HTTP Status: " << http_code << ", " << HTTP_STATUS_MESSAGES.at(http_code) << endl;
        html_logger.add_row("background-color:palegreen", row);
        return 0;
    case 202:
        cout << request_type << " accept request suceeded!" << endl;
        cout << "HTTP Status: " << http_code << ", " << HTTP_STATUS_MESSAGES.at(http_code) << endl;
        html_logger.add_row("background-color:palegreen", row);
        return 0;
    default:
        cout << request_type << " request failed!" << endl;
        cerr << "HTTP status code: " << http_code << ", " << HTTP_STATUS_MESSAGES.at(http_code) << endl;
        html_logger.add_row("background-color:tomato", row);
        return 1;
    }

    cerr << "Error not caught within responses!" << endl;
    return 1;
}

void clean_requests(CURL *curl, FILE *output_file)
{
    curl_easy_cleanup(curl);
    curl_global_cleanup();
    fclose(output_file);
}

int request_sender(FILE *output_file, CURL *curl, string request_type, string body)
{
    long http_code;
    CURLcode res;
    string res_string;

    http_code = 0;

    // GET/POST requests set up
    if (request_type == "GET")
    {
        res = curl_easy_perform(curl);
        curl_easy_getinfo(curl, CURLINFO_RESPONSE_CODE, &http_code);
        check_response(res, http_code, request_type, body);
        return http_code;
    }
    else if (request_type == "POST")
    {
        cout << body << endl;
        curl_easy_setopt(curl, CURLOPT_POSTFIELDS, body.c_str());

        // Setting headers
        headers = NULL;
        headers = curl_slist_append(headers, "Content-Type: application/json");
        curl_easy_setopt(curl, CURLOPT_HTTPHEADER, headers);

        // POST request performed
        res = curl_easy_perform(curl);

        curl_easy_getinfo(curl, CURLINFO_RESPONSE_CODE, &http_code);

        // Parse status code
        check_response(res, http_code, request_type, body);

        // Clean headers
        curl_slist_free_all(headers);

        return http_code;
    }
    else
    {
        cerr << "Invalid request type: " << request_type << endl;
        return http_code;
    }
    return http_code;
}

void initialise_requests(string url)
{
    curl = nullptr;
    output_file = nullptr;

    curl_global_init(CURL_GLOBAL_ALL);
    curl = curl_easy_init();
    if (!curl)
    {
        cerr << "curl_easy_init() failed" << endl;
    }

    // Setting output file. TODO: Modify this according to the mutation and fuzzer
    output_file = fopen("./src/fuzzing_responses/response.txt", "ab");
    if (!output_file)
    {
        cerr << "Failed to open output file" << endl;
        curl_easy_cleanup(curl);
        curl_global_cleanup();
    }

    curl_easy_setopt(curl, CURLOPT_URL, url.c_str());
    curl_easy_setopt(curl, CURLOPT_WRITEDATA, output_file);
    curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, write_callback);
}

int Django_Test_Driver(int energy, string url, string request_type, string input_file_path)
{
    int accumulated_iterations;
    bool testing_incomplete;
    list<json_seed> input_q;

    // create html logger file
    html_logger.create_file();
    vector<string> column_names = {"Time", "Request type", "Sent Contents", "HTTP Code"};
    html_logger.create_table_headings("background-color:lightgrey", column_names);

    // Check if the testing is complete, if so break out of while loop
    testing_incomplete = true;

    // Iterations of the given test, iterations 1 because energy starts at 1.
    accumulated_iterations = 0;

    // Initialise all curl requests with the url and output file
    initialise_requests(url);

    json json_input = parse_json(input_file_path);

    for (auto &el : json_input.items())
    {
        json_seed seed;
        seed.data = json_input;
        seed.key_to_mutate = el.key();
        seed.energy = 3;
        input_q.push_back(seed);
    }

    while (testing_incomplete)
    {

        json_seed cur_seed = input_q.front();
        input_q.pop_front();

        for (int i = 0; i < cur_seed.energy; i++)
        {
            // TODO: Mutate
            //  string value_to_mutate = cur_seed.data[cur_seed.key_to_mutate];
            //  string mutated_string = bit_flip(value_to_mutate);
            // TODO: Mutate ^^^
            // cur_seed.data[cur_seed.key_to_mutate] = mutated_string;
            cur_seed.data = mutate_requests(request_type, cur_seed.data);
            input_q.push_back(cur_seed);
            string json_body = cur_seed.data.dump();

            long http_code = request_sender(output_file, curl, request_type, json_body);

            // Check for interesting inputs
            // TODO: add output file path here too.
            string res_string = get_last_line();
            if (!is_interesting(res_string, http_code))
            {
                // Not interesting so remove new mutated input
                input_q.pop_back();
            }
        }

        accumulated_iterations++;

        // Change iterations here:
        if (accumulated_iterations > 10)
        {
            testing_incomplete = false;
        }
    }

    html_logger.close_file();

    // Dealloc all memory allocated to file, headers and curl
    clean_requests(curl, output_file);
    return 0;
}