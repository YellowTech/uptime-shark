import { createStore } from 'vuex'
import * as iots from 'io-ts'
import { withFallback } from "io-ts-types/lib/withFallback";
import { isRight } from 'fp-ts/Either'
import { fetchTimeout } from '@/util/utils';

const logEntry = iots.type({
    Failed: iots.boolean,
    Message: iots.string,
    Time: iots.number,
});

const monitor = iots.type({
    id: iots.string,
    name: iots.string,
    interval: withFallback(iots.number, 0),
    enabled: withFallback(iots.boolean, false),
    status: withFallback(iots.boolean, false),
    statusMessage: iots.string,
    inverted: withFallback(iots.boolean, false),
    mode: withFallback(iots.string, "undefined"),
    url: withFallback(iots.string, "undefined"),
    logs: withFallback(iots.array(logEntry), []),
});

type MonitorEdit = {
    id: string,
    name: string,
    interval: number,
    enabled: boolean,
    inverted: boolean,
    mode: string,
    url: string,
}

const monitors = iots.array(monitor)
// export interface Monitor = monitor._A
type Monitors = typeof monitors._A
type Monitor = typeof monitor._A


export type {
  Monitors,
  Monitor,
  MonitorEdit
}

export default createStore({
    state: {
        counter: 0,
        apiDomain: "",
        error: false,
        errorMessage: "",
        loaded: false,
        monitors: monitors._A,
        authenticated: false,
    },
    mutations: {
        increment(state) {
            state.counter++;
        },
    },
    actions: {
        increment(context) {
            context.commit("increment");
        },

        fetchData() {
            fetchTimeout(this.state.apiDomain + "/api/status", 3000, {
                credentials: "include",
            })
                .then((response) => response.text())
                .then((text) => {
                    const data = JSON.parse(text.replace("while(1);", ""));
                    const result = monitors.decode(data);
                    if (isRight(result)) {
                        // eslint-disable-next-line
                        this.state.monitors = result.right;
                        this.state.loaded = true;
                        this.state.error = false;
                    } else {
                        console.log(`Error when decoding!`);
                        console.log(data);
                        throw "ERROR: api data is incorrect";
                    }
                })
                .catch((error) => {
                    console.log("Error when fetching: " + error.name);
                    this.state.loaded = true;
                    this.state.error = true;
                    this.state.errorMessage = "Error reaching the server";

                    if (error.name === "AbortError") {
                        console.log("Timeout error");
                    } else if (error.name === "AbortError") {
                        console.log(`Error when decoding!`);
                    }
                });
        },

        checkLogin() {
            fetchTimeout(this.state.apiDomain + "/api/auth/status", 3000, {
                credentials: "include",
            })
                .then((response) => {
                    if (!response.ok) {
                        this.state.authenticated = false;
                    } else {
                        this.state.authenticated = true;
                    }
                })
                .catch((error) => {
                    console.log(
                        "Error when checking login status: " + error.name
                    );
                });
        },
    },
    modules: {},
    getters: {
        absCounter(state) {
            return Math.abs(state.counter);
        },
    },
});
