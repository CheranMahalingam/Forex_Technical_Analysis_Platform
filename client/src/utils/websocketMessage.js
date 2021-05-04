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

export const parseNewInferenceData = (parsedData, pastInferenceData) => {
    const symbol = Object.keys(parsedData)[0];
    let newInferenceData = [];
    if (!parsedData[symbol].inference) {
        return pastInferenceData;
    }
    const date = new Date();
    for (let i = 0; i < parsedData[symbol].inference.length; i++) {
        let newDate = new Date();
        newDate.setMinutes(date.getMinutes() + 15*(i + 1));
        newInferenceData.push({timestamp: newDate, inference: parsedData[symbol].inference[i]});
    }
    return newInferenceData;
}
