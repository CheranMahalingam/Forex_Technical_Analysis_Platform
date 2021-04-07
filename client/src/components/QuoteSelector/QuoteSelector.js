import React, { useState, useEffect } from "react";
import TextField from "@material-ui/core/TextField";
import Autocomplete from "@material-ui/lab/Autocomplete";
import Chip from "@material-ui/core/Chip";
import { makeStyles } from "@material-ui/core/styles";
import { parseISO } from "date-fns";
import Quote from "../Quote/Quote";
import { findMissingElement } from "../../utils/arrayUtils";
import { removeProperty } from "../../utils/parseObjectUtils";

var socket = new WebSocket("ws://localhost:8080/ws");

const useStyles = makeStyles(() => ({
  root: {
    "& .MuiOutlinedInput-root .MuiOutlinedInput-notchedOutline": {
      borderColor: "#FFFFFF",
    },
    "&:hover .MuiOutlinedInput-root .MuiOutlinedInput-notchedOutline": {
      borderColor: "#FFFFFF",
    },
    "& .MuiOutlinedInput-input": {
      color: "#FFFFFF",
    },
    "&:hover .MuiOutlinedInput-input": {
      color: "#FFFFFF",
    },
    "& .MuiOutlinedInput-root.Mui-focused .MuiOutlinedInput-input": {
      color: "#FFFFFF",
    },
    "& .MuiInputLabel-outlined": {
      color: "#FFFFFF",
    },
    "&:hover .MuiInputLabel-outlined": {
      color: "#FFFFFF",
    },
    "& .MuiAutocomplete-popupIndicator": {
      color: "#FFFFFF",
    },
  },
}));

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
  const [quoteHeight, setQuoteHeight] = useState({});
  const [selectedIndicators, setSelectedIndicators] = useState({});
  const [indicatorHeight, setIndicatorHeight] = useState({});

  const classes = useStyles();

  const height = 750;
  const chartHeight = height / 3;
  const macdHeight = (height * 11) / 72;
  const rsiHeight = (height * 11) / 72;
  const stoHeight = (height * 11) / 72;
  const spaceHeight = height / 54;
  const searchHeight = 150;
  const chartSpacingHeight = 120;

  const pushNewData = (newData) => {
    let parsedData = JSON.parse(newData);
    const symbol = Object.keys(parsedData)[0];
    let newPairData = parsedData[symbol];
    for (let i = 0; i < newPairData.length; i++) {
      newPairData[i].timestamp = parseISO(newPairData[i].timestamp);
    }
    let copyObj = { ...data };
    copyObj[symbol].push(...newPairData);
    console.log(data);
    return copyObj;
  };

  useEffect(() => {
    socket = new WebSocket("ws://localhost:8080/ws");

    socket.onopen = () => {
      console.log("websocket connected");
    };

    socket.onclose = () => {
      console.log("websocket disconnected");
    };

    socket.onmessage = (evt) => {
      console.log("onmessage");
      setData(pushNewData(evt.data));
    };

    socket.onerror = (err) => {
      console.log("encountered error:", err);
      socket.close();
    };

    return () => {
      console.log("component unmounted");
      socket.close();
    };
  }, []);

  const handleQuoteHeight = (symbol, heightChange) => {
    let newHeight = adjustHeight(symbol, heightChange);
    newHeight[symbol][1] = newHeight[symbol][1] + heightChange;
    setQuoteHeight(newHeight);
  };

  const adjustHeight = (symbol, heightChange) => {
    let newHeight = { ...quoteHeight };
    Object.keys(newHeight).forEach((currentSymbol) => {
      if (newHeight[currentSymbol][0] > newHeight[symbol][0]) {
        newHeight[currentSymbol][0] =
          newHeight[currentSymbol][0] + heightChange;
      }
    });
    return newHeight;
  };

  const handlePairChange = (_, selected) => {
    if (selected.length > listedPairs.length) {
      let subscribedSymbol = findMissingElement(selected, listedPairs);
      socket.send(
        JSON.stringify({ messageType: "subscribe", symbol: subscribedSymbol })
      );
      let newHeight;
      if (Object.keys(quoteHeight).length !== 0) {
        newHeight =
          Math.max(
            ...Object.values(quoteHeight).map((height) => height[0] + height[1])
          ) +
          chartHeight +
          chartSpacingHeight;
      } else {
        newHeight = searchHeight;
      }
      let newHeightObj = { ...quoteHeight, [subscribedSymbol]: [newHeight, 0] };
      console.log(newHeightObj);
      setQuoteHeight(newHeightObj);
      setSelectedIndicators({ ...selectedIndicators, [subscribedSymbol]: [] });
      setIndicatorHeight({
        ...indicatorHeight,
        [subscribedSymbol]: [chartHeight, 0, 0, 0],
      });
    } else {
      let unsubscribedSymbol = findMissingElement(listedPairs, selected);
      socket.send(
        JSON.stringify({
          messageType: "unsubscribe",
          symbol: unsubscribedSymbol,
        })
      );
      setData(removeProperty({ ...data }, unsubscribedSymbol));
      let newHeight = adjustHeight(
        unsubscribedSymbol,
        -chartHeight - chartSpacingHeight - quoteHeight[unsubscribedSymbol][1]
      );
      newHeight = removeProperty(newHeight, unsubscribedSymbol);
      setQuoteHeight(newHeight);
      const newIndicatorHeight = removeProperty(
        { ...indicatorHeight },
        unsubscribedSymbol
      );
      setIndicatorHeight(newIndicatorHeight);
      const newSelectedIndicators = removeProperty(
        { ...selectedIndicators },
        unsubscribedSymbol
      );
      setSelectedIndicators(newSelectedIndicators);
    }
    setListedPairs([...selected]);
  };

  const handleIndicatorsChange = (event, symbol) => {
    if (event.target.checked) {
      let newIndicator = event.target.name;
      let newHeight = indicatorHeight[symbol];
      if (newIndicator === "MACD") {
        newHeight[1] = Math.max(...newHeight) + macdHeight + spaceHeight;
        handleQuoteHeight(symbol, macdHeight + spaceHeight);
      } else if (newIndicator === "Relative Strength Index") {
        newHeight[2] = Math.max(...newHeight) + rsiHeight + spaceHeight;
        handleQuoteHeight(symbol, rsiHeight + spaceHeight);
      } else if (newIndicator === "Stochastic Oscillator") {
        newHeight[3] = Math.max(...newHeight) + stoHeight + spaceHeight;
        handleQuoteHeight(symbol, stoHeight + spaceHeight);
      }
      setIndicatorHeight({ ...indicatorHeight, [symbol]: [...newHeight] });
      setSelectedIndicators({
        ...selectedIndicators,
        [symbol]: [...selectedIndicators[symbol], newIndicator],
      });
    } else {
      let removedIndicator = event.target.name;
      let newHeight = indicatorHeight[symbol];
      if (removedIndicator === "MACD") {
        newHeight = newHeight.map((height) =>
          height > newHeight[1] + spaceHeight
            ? height - macdHeight - spaceHeight
            : height
        );
        handleQuoteHeight(symbol, -macdHeight - spaceHeight);
        newHeight[1] = 0;
      } else if (removedIndicator === "Relative Strength Index") {
        newHeight = newHeight.map((height) =>
          height > newHeight[2] + spaceHeight
            ? height - rsiHeight - spaceHeight
            : height
        );
        handleQuoteHeight(symbol, -rsiHeight - spaceHeight);
        newHeight[2] = 0;
      } else if (removedIndicator === "Stochastic Oscillator") {
        newHeight = newHeight.map((height) =>
          height > newHeight[3] + spaceHeight
            ? height - stoHeight - spaceHeight
            : height
        );
        handleQuoteHeight(symbol, -stoHeight - spaceHeight);
        newHeight[3] = 0;
      }
      setIndicatorHeight({ ...indicatorHeight, [symbol]: [...newHeight] });
      setSelectedIndicators({
        ...selectedIndicators,
        [symbol]: selectedIndicators[symbol].filter(
          (name) => name !== removedIndicator
        ),
      });
    }
  };

  return (
    <React.Fragment>
      <Autocomplete
        multiple
        id="quote-tags"
        limitTags={5}
        classes={classes}
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
            classes={classes}
          />
        )}
        style={{
          color: "#FFFFFF",
          position: "relative",
          marginLeft: "20%",
          marginRight: "20%",
          marginTop: "2%",
          marginBottom: "1%",
        }}
      />
      {listedPairs &&
        listedPairs.map((pair, index) => {
          return (
            <Quote
              key={index}
              data={data[pair]}
              pair={pair}
              quoteHeight={quoteHeight}
              handleQuoteHeight={handleQuoteHeight}
              handleIndicatorsChange={handleIndicatorsChange}
              selectedIndicators={selectedIndicators[pair]}
              indicatorHeight={indicatorHeight[pair]}
              macdHeight={macdHeight}
              rsiHeight={rsiHeight}
              stoHeight={stoHeight}
            />
          );
        })}
    </React.Fragment>
  );
}

const currencyPairs = [
  "EURUSD",
  "GBPUSD",
  "USDCAD",
  "AUDUSD",
  "USDJPY",
  "USDCHF",
  "NZDJPY",
];

export default QuoteSelector;
