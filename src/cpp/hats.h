#include <vector>
#include <cstdint>

typedef std::vector<std::vector<float>> ha_array;

typedef std::vector<uint8_t> buffer;

typedef  std::vector<std::vector<int>> cidx_map;

struct event {
    int x, y, p;
    float t;
};

typedef std::vector<event> stream;

struct params {
    int R, K, width, height;
    float tau, delta, time_window;
    cidx_map *idx;
};

struct hats {
    int n_cells, width, height, dim;
    std::vector<ha_array> data [2];
    std::vector<stream> mc [2];
    cidx_map *idx;
    std::vector<int> evc [2];
};

ha_array make_ha_array(int width, int height) {
    ha_array out = {};

    for (int i = 0; i < width; i++) {
        std::vector<float> e;
        out.push_back(e);
        for (int j = 0; j < height; j++)
            out[i].push_back(0);
    }
    return out;
}

cidx_map make_idx(int width, int height, int K) {
    cidx_map out;
    int cell_width = width / K;
    for (int i = 0; i < width; i++) {
        std::vector<int> e;
        out.push_back(e);
        for (int j=0; j < height; i++) {
            int p_row = i / K;
            int p_col = j / K;
            out[i][j] = p_row*cell_width + p_col;
        }
    }
    return out;
}

hats make_hats(params *p, cidx_map *idx) { // int n_cells, int width, int height, int dim, indexarr *idx) {
    hats out = {};
    out.width = (*p).width;
    out.height = (*p).height;
    out.n_cells = (out.width/(*p).K)*(out.height/(*p).K);
    out.dim = 2*out.n_cells*(2*(*p).R + 1)*(2*(*p).R + 1);
    out.idx = idx;

    for (int i = 0; i < 2; i++) {
        std::vector<ha_array> data;
        std::vector<stream> mc;
        out.data[i] = data;
        out.mc[i] = mc;

        for (int j = 0; j < out.n_cells; j++) // needed ?
            out.data[i].push_back(make_ha_array(out.width, out.height));
    }
    return out;
}