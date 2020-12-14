let stocks = ((api_base_url, api_path) => {
    let api_symbols_url = new URL(api_path, api_base_url);
    let ticker_symbols = fetch(api_symbols_url.toString()).then(res => res.json());
    
    let getTickerSymbols = () => ticker_symbols;

    return {
        getTickerSymbols: getTickerSymbols,
    };
})("https://5bdexayy1k.execute-api.us-west-2.amazonaws.com","dev/symbols");

export default stocks;