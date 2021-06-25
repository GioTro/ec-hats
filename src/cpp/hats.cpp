#include <hats.h>
#include <stdlib.h>
#include <iostream>
#include <vector>
#include <bits/stdc++.h>

using namespace std;

int time_surface(event e, params *p, stream *mce, vector<vector<float>> *hst) {
    for (int i = 0; i < (*mce).size(); i++) {
        float dt = (e.t - (*mce)[i].t);
        int shifted_y = (*mce)[i].y - (e.y - (*p).R);
        int shifted_x = (*mce)[i].x - (e.x - (*p).R);
        (*hst)[shifted_y][shifted_x] += exp(-dt/(*p).tau);
    }
}

void process(event e, params *p, hats *ds) {
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

    time_surface(e, p, &(*ds).mc[e.p][idx], &(*ds).data[e.p][idx]);
}

int process_all() {
    return 0;
}