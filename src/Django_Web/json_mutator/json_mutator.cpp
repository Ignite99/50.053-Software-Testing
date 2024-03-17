#include <curl/curl.h>
#include "nlohmann/json.hpp"

#include <iostream>
#include <fstream>
#include <string>
#include <sstream>
#include <random>

using namespace std;
using json = nlohmann::json;

/* ============================================================================ */
/* ============================= HELPER FUNCTIONS ============================= */
/* ============================================================================ */
string generate_random_string(int length)
{
    // Define a pool of characters from which to select
    const string charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789";

    // Create a random number engine and seed it with a random device
    random_device rd;
    mt19937 gen(rd());

    // Create a distribution to select characters from the charset
    uniform_int_distribution<> dis(0, charset.length() - 1);

    // Generate the random string
    string random_string;
    for (int i = 0; i < length; ++i)
    {
        random_string += charset[dis(gen)];
    }

    return random_string;
}

bool string_convertible_to_int(const string &str)
{
    stringstream ss(str);
    int num;
    ss >> num;
    return !ss.fail() && ss.eof();
}

void change_value_type(json &data, const string key)
{
    // random device generators
    random_device rd;
    mt19937 gen(rd());
    uniform_int_distribution<> dis(0, 1);
    int random_number;

    // check if key exists
    if (data.find(key) != data.end())
    {
        auto value = data[key];
        json::value_t val_type = data[key].type();
        cout << "Update at key: " << key << endl;

        if (val_type == json::value_t::number_unsigned)
        {
            cout << "Value is unsigned integer" << endl;
            // convert to float or string
            random_number = dis(gen);
            json::number_unsigned_t unsigned_val = value.get<json::number_unsigned_t>();
            if (random_number == 0)
            {
                float float_res = static_cast<float>(unsigned_val);
                cout << "Convert to float" << endl;
                data[key] = float_res;
            }
            else if (random_number == 1)
            {
                string string_res = to_string(unsigned_val);
                cout << "Convert to string" << endl;
                data[key] = string_res;
            }
        }
        else if (val_type == json::value_t::number_float)
        {
            cout << "Value is float" << endl;
        }
        else if (val_type == json::value_t::string)
        {
            cout << "Value is string" << endl;
            json::string_t string_val = value.get<json::string_t>();
            if (string_convertible_to_int(value))
            {
                int int_val = stoi(string_val);
                cout << "Convert to int" << endl;
                data[key] = int_val;
            }
        }
        else
        {
            cout << "No valid type found for this value" << endl;
        }
    }
}

/* ================================================================================= */
/* ============================= MAIN MUTATOR FUNCTION ============================= */
/* ================================================================================= */
json mutate_requests(string request_type, json data)
{
    // initialize variables
    string new_key;
    string new_val;
    int index;
    auto it = data.begin();
    // create random device generator
    random_device rd;
    mt19937 gen(rd());
    uniform_int_distribution<> dis(1, 3);

    if (request_type == "POST")
    {
        // randomize options
        int random_number = dis(gen);
        if (random_number == 2 && data.empty())
        {
            // add a new JSON field
            random_number = 1;
        }

        switch (random_number)
        {
        case 1:
            // 1. Add new JSON field
            cout << "Mutation 1: Add new JSON field" << endl;
            new_key = generate_random_string(5);
            new_val = generate_random_string(5);
            data[new_key] = new_val;
            break;
        case 2:
            // 2. Remove existing JSON field
            cout << "Mutation 2: Remove existing JSON field" << endl;
            index = dis(gen);
            it = next(data.begin(), index);
            data.erase(it);
            break;
        case 3:
            // 3. Keep all fields as it is but change value type
            cout << "Mutation 3: Change type of existing JSON field value" << endl;
            index = dis(gen);
            it = next(data.begin(), index);
            change_value_type(data, it.key());
            break;
        }
    }
    else if (request_type == "GET")
    {
        // TODO: handle GET request JSON mutations
    }

    return data;
}

int main()
{
    string input_file_path = "../template_inputs/add_product.json";
    ifstream input_file(input_file_path);
    if (!input_file.is_open())
    {
        cerr << "Error: Could not open file '" << input_file_path << "'." << endl;
        return 1;
    }

    json data;
    input_file >> data;
    input_file.close();

    json new_data = mutate_requests("POST", data);

    cout << "JSON data: " << endl;
    cout << new_data.dump(4) << endl;

    return 0;
}