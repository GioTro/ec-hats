#include <vector>
#include <cstdint>

using namespace std;

typedef vector<vector<float>> ha_array;

typedef vector<uint8_t> buffer;

typedef vector<vector<int>> indexarr;

struct event {
    int x, y, p;
    float t;
};

typedef vector<event> stream;

struct params {
    int R, K, width, height;
    float tau, delta, time_window;
};

struct hats {
    int n_cells, width, height, dim;
    vector<ha_array> data [2];
    vector<stream> mc [2];
    indexarr *idx;
    vector<int> evc [2];
};

ha_array make_ha_array(int width, int height) {
    ha_array out = {};

    for (int i = 0; i < width; i++) {
        vector<float> e;
        out.push_back(e);
        for (int j = 0; j < height; j++)
            out[i].push_back(0);
    }
    return out;
}

vector<vector<int>> make_idx(int width, int height, int K) {
    vector<vector<int>> out;
    int cell_width = width / K;
    for (int i = 0; i < width; i++) {
        vector<int> e;
        out.push_back(e);
        for (int j=0; j < height; i++) {
            int p_row = i / K;
            int p_col = j / K;
            out[i][j] = p_row*cell_width + p_col;
        }
    }
    return out;
}

hats make_hats(int n_cells, int width, int height, int dim, indexarr *idx) {
    hats out = {};
    out.n_cells = n_cells;
    out.width = width;
    out.height = height;
    out.dim = dim;
    out.idx = idx;

    for (int i = 0; i < 2; i++) {
        vector<ha_array> data;
        vector<stream> mc;
        out.data[i] = data;
        out.mc[i] = mc;

        for (int j = 0; j < n_cells; j++) // needed ?
            out.data[i].push_back(make_ha_array(width, height));
    }
    return out;
}