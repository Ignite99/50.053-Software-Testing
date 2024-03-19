#include <curl/curl.h>
#include "nlohmann/json.hpp"

#include <iostream>
#include <fstream>
#include <string>
#include <sstream>
#include <random>

using namespace std; // causes ambiguity for cout
using json = nlohmann::json;

/* ============================================================================ */
/* ============================= HELPER FUNCTIONS ============================= */
/* ============================================================================ */
string generate_random_string(int length)
{
    const string charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789";
    random_device rd;
    mt19937 gen(rd());
    uniform_int_distribution<> dis(0, charset.length() - 1);

    string random_string;
    for (int i = 0; i < length; ++i)
    {
        random_string += charset[dis(gen)];
    }
    return random_string;
}

// checks if a string is convertible to integer
bool string_convertible_to_int(const string &str)
{
    stringstream ss(str);
    int num;
    ss >> num;
    return !ss.fail() && ss.eof();
}

// return either true or false
bool choose_boolean_output()
{
    random_device rd;
    mt19937 gen(rd());
    uniform_int_distribution<> dis(0, 1);
    int random_number = dis(gen);
    if (random_number)
        return false;
    else
        return true;
}

// changes one value type to another type
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

        // 1. UNSIGNED INT or INTEGER
        if (val_type == json::value_t::number_unsigned || val_type == json::value_t::number_integer)
        {
            cout << "Value is unsigned integer" << endl;
            unsigned int num_val;
            // convert INT to FLOAT, STRING, BOOL
            random_number = dis(gen);
            if (val_type == json::value_t::number_unsigned)
                num_val = value.get<json::number_unsigned_t>();
            else
                num_val = value.get<json::number_integer_t>();
            if (random_number == 0)
            {
                float float_res = static_cast<float>(num_val);
                cout << "Convert to float" << endl;
                data[key] = float_res;
            }
            else if (random_number == 1)
            {
                string string_res = to_string(num_val);
                cout << "Convert to string" << endl;
                data[key] = string_res;
            }
        }
        // 2. FLOAT
        else if (val_type == json::value_t::number_float)
        {
            cout << "Value is float" << endl;
            // convert FLOAT to INT
            cout << "Convert to int" << endl;
            data[key] = static_cast<int>(value);
        }
        // 3. STRING
        else if (val_type == json::value_t::string)
        {
            cout << "Value is string" << endl;
            // convert STRING to INT, BOOL
            json::string_t string_val = value.get<json::string_t>();
            if (string_convertible_to_int(value))
            {
                int int_val = stoi(string_val);
                cout << "Convert to int" << endl;
                data[key] = int_val;
            }
            // TODO: update strings that have mix of characters and numbers
        }
        // 4. BOOLEAN
        else if (val_type == json::value_t::boolean)
        {
            cout << "Value is boolean" << endl;
            cout << "Convert to string" << endl;
            if (data[key])
                data[key] = "true";
            else
                data[key] = "false";
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
    // initialize variables for switch block
    string new_key;
    string new_val;
    int index;
    auto it = data.begin();
    float random_float;

    // create random device generator
    random_device rd;
    mt19937 gen(rd());
    uniform_int_distribution<> mutate_dis(1, 3);     // 3 different mutations
    uniform_int_distribution<> random_gen_dis(1, 4); // 4 different value types
    uniform_int_distribution<> index_dis(0, data.size() - 1);
    uniform_int_distribution<> num_dis(0, 100);

    if (request_type == "POST")
    {
        // randomize options
        int random_number = mutate_dis(gen);
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
            // randomize values generated
            random_number = random_gen_dis(gen);
            switch (random_number)
            {
            case 1:
                new_val = generate_random_string(5);
                break;
            case 2:
                new_val = num_dis(gen);
                break;
            case 3:
                random_float = num_dis(gen);
                new_val = random_float;
            case 4:
                new_val = choose_boolean_output();
            }
            data.emplace(new_key, new_val);
            break;
        case 2:
            // 2. Remove existing JSON field
            index = index_dis(gen);
            cout << "Mutation 2: Remove existing JSON field at index " << index << endl;
            it = next(data.begin(), index);
            data.erase(it);
            break;
        case 3:
            // 3. Keep all fields as it is but change value type
            cout << "Mutation 3: Change type of existing JSON field value" << endl;
            index = index_dis(gen);
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

int main(int argc, char *argv[])
{
    // command line arguments
    int rounds = 10;
    if (argc > 1)
    {
        rounds = stoi(argv[1]);
    }
    cout << "Number of mutation rounds: " << rounds << endl;

    // JSON input file
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

    // loop through "rounds" and mutate in each iteration
    json output_data = data;
    for (int i = 0; i < rounds; i++)
    {
        output_data = mutate_requests("POST", output_data);
        cout << output_data.dump(4) << endl;
    }

    return 0;
}