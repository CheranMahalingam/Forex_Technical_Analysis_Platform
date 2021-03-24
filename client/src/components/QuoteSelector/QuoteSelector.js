import React, { useState, useEffect } from "react";
import TextField from "@material-ui/core/TextField";
import Autocomplete from "@material-ui/lab/Autocomplete";
import Chip from "@material-ui/core/Chip";
import { parseISO } from "date-fns";
import Quote from "../Quote/Quote";
import nice from "../Data";
import fakeData from "../../assets/fake_data.json";
import "./QuoteSelector.css";

function QuoteSelector() {
  const [listedPairs, setListedPairs] = useState([]);

  useEffect(() => {
    const socket = new WebSocket("ws://localhost:8080/ws");
    socket.onmessage = (evt) => console.log(evt.data);
  }, []);

  const handlePairChange = (event, new_value) => {
    setListedPairs([...new_value]);
  };

  const cool = fakeData.map((data) => {
    const date = parseISO(data.date);
    return {
      date: date,
      high: data.high,
      low: data.low,
      open: data.open,
      close: data.close,
    };
  });

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
        {listedPairs &&
          listedPairs.map((pair, index) => {
            return <Quote key={index} data={nice} pair={pair} />;
          })}
      </div>
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
