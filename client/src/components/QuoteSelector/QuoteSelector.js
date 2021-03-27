import React, { useState, useEffect } from "react";
import TextField from "@material-ui/core/TextField";
import Autocomplete from "@material-ui/lab/Autocomplete";
import Chip from "@material-ui/core/Chip";
import { parseISO } from "date-fns";
import Quote from "../Quote/Quote";
import "./QuoteSelector.css";

function QuoteSelector() {
  const [listedPairs, setListedPairs] = useState([]);
  const [data, setData] = useState([]);

  useEffect(() => {
    const socket = new WebSocket("ws://localhost:8080/ws");
    socket.onmessage = (evt) => {
      let parsedData = JSON.parse(evt.data);
      console.log(parsedData);
      for (let i = 0; i < parsedData.length; i++) {
        parsedData[i].timestamp = parseISO(parsedData[i].timestamp);
      }
      console.log([...data, ...parsedData]);
      setData((data) => [...data, ...parsedData]);
    };
  }, []);

  const handlePairChange = (event, new_value) => {
    setListedPairs([...new_value]);
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
          return <Quote key={index} data={data} pair={pair} />;
        })}
    </React.Fragment>
  );
}

const currencyPairs = [
  "EUR/USD",
  "GBP/USD",
  "USD/CAD",
  "USD/JPY",
  "AUD/USD",
  "USD/CHF",
  "NZD/JPY",
];

export default QuoteSelector;
