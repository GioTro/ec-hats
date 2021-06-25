#include <hats.h>
#include <stdlib.h>
#include <iostream>
#include <vector>
#include <bits/stdc++.h>

using namespace std;

int time_surface(event e, params *p, stream *mce, vector<vector<float>> *hst) {
    for (int i = 0; i < (*mce).data.size(); i++) {
        float dt = (e.t - (*mce).data[i].t);

        int shifted_y = (*mce).data[i].y - (e.y - (*p).R);
        int shifted_x = (*mce).data[i].x - (e.x - (*p).R);
        (*hst)[shifted_y][shifted_x] += exp(-dt/(*p).tau);
    }
}

void process(event e, params *p, hats *ds) {
    int idx = (*(*ds).idx)[e.x][e.y];
    auto data = &((*ds).mc[e.p][idx].data);
    (*ds).evc[e.p][idx]++;
    (*data).push_back(e);
    float bp = e.t - (*p).time_window;
    int i = 0;

    while ((*data)[i].t < bp && i < (*data).size()) {
        i++;
    }

    (*data) = vector<event>((*data).begin() + i, (*data).end());
    time_surface(e, p, &(*ds).mc[e.p][idx], &(*ds).data[e.p][idx].data);
}

int process_all() {
    return 0;
}