export const fetchTimeout = (url: RequestInfo, ms: number, init: RequestInit|undefined=undefined):Promise<Response> => {
    const controller = new AbortController();
    const fetchRequestInit: RequestInit = (init ? init : {})
    fetchRequestInit.signal = controller.signal
    const promise = fetch(url, fetchRequestInit);
    init?.signal?.addEventListener("abort", () => controller.abort());
    const timeout = setTimeout(() => controller.abort(), ms);
    return promise.finally(() => clearTimeout(timeout));
};
