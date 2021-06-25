#include <hats.h>
#include <stdlib.h>
#include <iostream>
#include <vector>
#include <bits/stdc++.h>

std::vector<float> normalize(hats *ds) {
    // Normalize and flatten
    std::vector<float> out;
    int idx = 0;

    for (int i = 0; i < sizeof((*ds).data); i++ )
        for (int j = 0; j < (*ds).data[i].size(); j++)
            for (int z = 0; z < (*ds).data[i][j].size(); z++)
                for (int k = 0; k < (*ds).data[i][j][z].size(); k++) {
                    out[idx] = (*ds).data[i][j][z][k] / ((float) (*ds).evc[i][j] + FLT_MIN);
                    idx++;
    }
    return out;
}

int time_surface(event e, params *p, stream *mce, ha_array *hst) {
    for (int i = 0; i < (*mce).size(); i++) {
        float dt = (e.t - (*mce)[i].t);
        int shifted_y = (*mce)[i].y - (e.y - (*p).R);
        int shifted_x = (*mce)[i].x - (e.x - (*p).R);
        (*hst)[shifted_y][shifted_x] += exp(-dt/(*p).tau);
    }
}

void process_event(event e, params *p, hats *ds) {
    int idx = (*(*ds).idx)[e.x][e.y];
    stream *mce = &(*ds).mc[e.p][idx];
    (*ds).evc[e.p][idx]++;
    (*mce).push_back(e);
    float bp = e.t - (*p).time_window;

    stream h;
    for (int i = 0; i < (*mce).size(); i++)
        if ((*mce)[i].t >= bp)
            h.push_back((*mce)[i]);

    (*mce) = h;

    time_surface(e, p, mce, &(*ds).data[e.p][idx]);
}

std::vector<float> process(stream es, params *p) {
    hats ds = make_hats(p, (*p).idx);

    for (auto e = begin(es); e < end(es); e++)
        process_event(*e, p, &ds);

    return normalize(&ds);
}