const { parseISO } = require("date-fns");

export const parseNewSymbolData = (parsedData, chartData) => {
  const symbol = Object.keys(parsedData)[0];
  let newPairData = parsedData[symbol];
  if (!newPairData) {
    return chartData;
  }
  for (let i = 0; i < newPairData.length; i++) {
    newPairData[i].timestamp = parseISO(newPairData[i].timestamp);
  }
  let copyObj = { ...chartData };
  copyObj[symbol].push(...newPairData);
  console.log(copyObj);
  return copyObj;
};

export const parseNewInferenceData = (
  parsedData,
  pastInferenceData,
  previous
) => {
  const symbol = Object.keys(parsedData)[0];
  let newInferenceData = [];
  if (!parsedData[symbol].inference) {
    return pastInferenceData;
  }
  const date = new Date();
  let slope = calculate_slope(parsedData[symbol].inference[0], previous, 15);
  for (let i = 0; i < parsedData[symbol].inference.length * 15; i++) {
    const index = Math.floor(i / 15);
    if (i % 15 === 0 && i !== 0) {
      slope = calculate_slope(
        parsedData[symbol].inference[index],
        parsedData[symbol].inference[index - 1],
        15
      );
    }
    const slopeSum = (i % 15) * slope;
    let newDate = new Date();
    newDate.setMinutes(date.getMinutes() + i + 1);
    if (index === 0) {
      newInferenceData.push({
        timestamp: newDate,
        inference: previous + slopeSum,
      });
    } else {
      newInferenceData.push({
        timestamp: newDate,
        inference: parsedData[symbol].inference[index - 1] + slopeSum,
      });
    }
  }
  return newInferenceData;
};

const calculate_slope = (y2, y1, xChange) => {
  if (xChange === 0) {
    return null;
  }
  return (y2 - y1) / xChange;
};
