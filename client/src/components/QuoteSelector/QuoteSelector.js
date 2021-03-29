import React, { useState, useEffect } from "react";
import TextField from "@material-ui/core/TextField";
import Autocomplete from "@material-ui/lab/Autocomplete";
import Chip from "@material-ui/core/Chip";
import { parseISO } from "date-fns";
import Quote from "../Quote/Quote";
import "./QuoteSelector.css";

const socket = new WebSocket("ws://localhost:8080/ws");

function QuoteSelector() {
  const [listedPairs, setListedPairs] = useState([]);
  const [data, setData] = useState({
    EURUSD: [],
    GBPUSD: [],
    USDCAD: [],
    USDJPY: [],
    AUDUSD: [],
    USDCHF: [],
    NZDJPY: [],
  });

  const pushNewData = (newData) => {
    let parsedData = JSON.parse(newData);
    const symbol = Object.keys(parsedData)[0];
    let newPairData = parsedData[symbol];
    for (let i = 0; i < newPairData.length; i++) {
      newPairData[i].timestamp = parseISO(newPairData[i].timestamp);
    }
    let copyData = { ...data };
    copyData[symbol].push(...newPairData);
    console.log(copyData);
    return copyData;
  };

  useEffect(() => {
    socket.onopen = () => {
      console.log("websocket connected");
    };

    socket.onclose = () => {
      console.log("websocket disconnected");
    };

    socket.onmessage = (evt) => {
      console.log("onmessage");
      setData(() => pushNewData(evt.data));
    };

    socket.onerror = (err) => {
      console.log("encountered error:", err);
      socket.close();
    };

    return () => {
      socket.close();
    };
  }, []);

  const handlePairChange = (_, selected) => {
    if (selected.length > listedPairs.length) {
      let subscribedSymbol = findNewSymbol(selected, listedPairs);
      console.log(subscribedSymbol);
      socket.send(
        JSON.stringify({ messageType: "subscribe", symbol: subscribedSymbol })
      );
    } else {
      let unsubscribedSymbol = findRemovedSymbol(selected, listedPairs);
      console.log(unsubscribedSymbol);
      socket.send(
        JSON.stringify({
          messageType: "unsubscribe",
          symbol: unsubscribedSymbol,
        })
      );
    }
    setListedPairs([...selected]);
  };

  const findNewSymbol = (selectedSymbols, subscribedSymbols) => {
    if (subscribedSymbols.length === 0) {
      return selectedSymbols[0];
    }
    for (let i = 0; i < selectedSymbols.length; i++) {
      let currentSymbol = selectedSymbols[i];
      for (let j = 0; j < subscribedSymbols.length; j++) {
        if (currentSymbol === subscribedSymbols[j]) {
          break;
        } else if (j === subscribedSymbols.length - 1) {
          return currentSymbol;
        }
      }
    }
    return null;
  };

  const findRemovedSymbol = (selectedSymbols, subscribedSymbols) => {
    if (selectedSymbols.length == 0) {
      return subscribedSymbols[0];
    }
    for (let i = 0; i < subscribedSymbols.length; i++) {
      let currentSymbol = subscribedSymbols[i];
      for (let j = 0; j < selectedSymbols.length; j++) {
        if (currentSymbol === selectedSymbols[j]) {
          break;
        } else if (j === selectedSymbols.length - 1) {
          return currentSymbol;
        }
      }
    }
    return null;
  };

  return (
    <React.Fragment>
      <div className="select-quote">
        <Autocomplete
          multiple
          id="quote-tags"
          limitTags={5}
          options={currencyPairs}
          filterSelectedOptions
          value={listedPairs}
          onChange={handlePairChange}
          renderTags={(tagValue, getTagProps) =>
            tagValue.map((option, index) => (
              <Chip label={option} {...getTagProps({ index })} />
            ))
          }
          renderInput={(params) => (
            <TextField
              {...params}
              variant="outlined"
              label="Currency Pairs"
              placeholder="Select a pair"
            />
          )}
        />
      </div>
      {listedPairs &&
        listedPairs.map((pair, index) => {
          return <Quote key={index} data={data[pair]} pair={pair} />;
        })}
    </React.Fragment>
  );
}

const currencyPairs = [
  "EURUSD",
  "GBPUSD",
  "USDCAD",
  "USDJPY",
  "AUDUSD",
  "USDCHF",
  "NZDJPY",
];

export default QuoteSelector;
