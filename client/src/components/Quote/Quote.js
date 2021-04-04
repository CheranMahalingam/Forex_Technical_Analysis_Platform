import React, { useState } from "react";
import { format } from "d3-format";
import { timeFormat } from "d3-time-format";
import Button from "@material-ui/core/Button";
import AddIcon from "@material-ui/icons/Add";
import { Chart, ChartCanvas } from "@react-financial-charts/core";
import { XAxis, YAxis } from "@react-financial-charts/axes";
import { discontinuousTimeScaleProviderBuilder } from "@react-financial-charts/scales";
import { CandlestickSeries } from "@react-financial-charts/series";
import { withDeviceRatio, withSize } from "@react-financial-charts/utils";
import {
  ema,
  bollingerBand,
  rsi,
  macd,
  stochasticOscillator,
} from "@react-financial-charts/indicators";
import {
  CrossHairCursor,
  MouseCoordinateY,
  MouseCoordinateX,
  mouseBasedZoomAnchor,
  ZoomButtons,
  EdgeIndicator,
  OHLCTooltip,
  Label,
  last,
} from "react-financial-charts";
import Rsi from "../Technical_Indicators/Rsi";
import Macd from "../Technical_Indicators/Macd";
import Stochastic from "../Technical_Indicators/Stochastic";
import BollingerBand from "../Technical_Indicators/Bollinger_Band";
import Ema from "../Technical_Indicators/Ema";
import SelectIndicatorModal from "../SelectIndicatorModal/SelectIndicatorModal";

function Quote(props) {
  const [openIndicators, setOpenIndicators] = useState(false);
  const [selectedIndicators, setSelectedIndicators] = useState([]);

  const ema12 = ema()
    .id(1)
    .options({ windowSize: 12 })
    .merge((d, c) => {
      d.ema12 = c;
    })
    .accessor((d) => d.ema12);
  const ema26 = ema()
    .id(2)
    .options({ windowSize: 26 })
    .merge((d, c) => {
      d.ema26 = c;
    })
    .accessor((d) => d.ema26);
  const bollingerCalculator = bollingerBand()
    .options({ windowSize: 20, movingAverageType: "sma" })
    .merge((d, c) => {
      d.bb = c;
    })
    .accessor((d) => d.bb);
  const rsiCalculator = rsi()
    .options({ windowSize: 14 })
    .merge((d, c) => {
      d.rsi = c;
    })
    .accessor((d) => d.rsi);
  const stoCalculator = stochasticOscillator()
    .options({ windowSize: 14, kWindowSize: 1, dWindowSize: 3 })
    .merge((d, c) => {
      d.fastSTO = c;
    })
    .accessor((d) => d.fastSTO);
  const macdCalculator = macd()
    .options({ fast: 12, signal: 9, slow: 26 })
    .merge((d, c) => {
      d.macd = c;
    })
    .accessor((d) => d.macd);
  const macdAppearance = {
    fillStyle: {
      divergence: "#4682B4",
    },
    strokeStyle: {
      macd: "#0093FF",
      signal: "#D84315",
      zero: "rgba(0, 0, 0, 0.3)",
    },
  };
  const candlesAppearance = {
    wickStroke: function fill(d) {
      return d.close < d.open ? "rgba(255, 7, 58, 1)" : "rgba(57, 255, 20, 1)";
    },
    fill: "#000000",
    stroke: function outline(d) {
      return d.close < d.open ? "rgba(255, 7, 58, 1)" : "rgba(57, 255, 20, 1)";
    },
    candleStrokeWidth: 1,
    width: 3,
  };

  const { type, height, width, ratio, data: initialData } = props;
  const chartHeight = height / 3;
  const macdHeight = (height * 12) / 72;
  const rsiHeight = (height * 12) / 72;
  const stoHeight = (height * 12) / 72;
  const spaceHeight = (1 * height) / 72;
  const calculatedData = macdCalculator(
    stoCalculator(rsiCalculator(bollingerCalculator(ema26(ema12(initialData)))))
  );
  const xScaleProvider = discontinuousTimeScaleProviderBuilder().inputDateAccessor(
    (d) => d.timestamp
  );
  const { data, xScale, xAccessor, displayXAccessor } = xScaleProvider(
    calculatedData
  );
  const max = xAccessor(last(data));
  const min = xAccessor(data.length > 200 ? data[data.length - 200] : data[0]);
  const xExtents = [min, max + 1];
  const pricesDisplayFormat = format(".5f");
  const changeDisplayFormat = format("+.5f");
  const margin = { left: 20, right: 80, top: 60, bottom: 30 };
  const pairLabel =
    currencyPairAlias[props.pair.slice(0, 3)] +
    " / " +
    currencyPairAlias[props.pair.slice(3)] +
    " \u0387 1 \u0387 OANDA";

  const handleColour = (d) => {
    if (d) {
      if (d.close === d.open) {
        return "#FFFFFF";
      }
      return d.close < d.open ? "rgba(255, 7, 58, 1)" : "rgba(57, 255, 20, 1)";
    } else {
      return "#FFFFFF";
    }
  };

  const handleIndicatorsOpen = () => {
    setOpenIndicators(true);
  };

  const handleIndicatorsClose = () => {
    setOpenIndicators(false);
  };

  const handleIndicatorsChange = (_, selected) => {
    console.log(selected);
    setSelectedIndicators([...selected]);
  };

  return (
    <React.Fragment>
      <Button
        variant="outlined"
        color="primary"
        style={{
          position: "absolute",
          left: width - 300,
          top: 20,
          zIndex: 7,
          color: "#FFFFFF",
          borderColor: "#FFFFFF",
        }}
        onClick={handleIndicatorsOpen}
        startIcon={<AddIcon />}
        size="small"
      >
        Indicators
      </Button>
      <SelectIndicatorModal
        openIndicators={openIndicators}
        handleIndicatorsClose={handleIndicatorsClose}
        selectedIndicators={selectedIndicators}
        handleIndicatorsChange={handleIndicatorsChange}
      />
      <ChartCanvas
        height={height}
        ratio={ratio}
        width={width}
        margin={margin}
        type={type}
        data={data}
        xAccessor={xAccessor}
        xScale={xScale}
        xExtents={xExtents}
        displayXAccessor={displayXAccessor}
        zoomAnchor={mouseBasedZoomAnchor}
      >
        <Label x={margin.left + 100} y={-15} fontSize={18} text={pairLabel} />
        <Chart
          id={1}
          yExtents={
            ((d) => [d.high, d.low],
            ema12.accessor(),
            ema26.accessor(),
            bollingerCalculator.accessor())
          }
          height={chartHeight}
        >
          {false ? (
            <XAxis
              strokeStyle="#FFFFFF"
              strokeDasharray="Solid"
              axisAt="bottom"
              orient="bottom"
              ticks={15}
              tickLabelFill="#FFFFFF"
              gridLinesStrokeDasharray="Solid"
            />
          ) : null}
          <YAxis
            showGridLines
            gridLinesStrokeWidth={0.5}
            strokeStyle="#FFFFFF"
            axisAt="right"
            orient="right"
            ticks={10}
            tickFormat={pricesDisplayFormat}
            tickLabelFill="#FFFFFF"
            gridLinesStrokeDasharray="Solid"
          />
          {selectedIndicators.includes("Bollinger Bands") ? (
            <BollingerBand bollingerCalculator={bollingerCalculator} />
          ) : null}
          {selectedIndicators.includes("Exponential Moving Average") ? (
            <Ema ema={[ema12, ema26]} />
          ) : null}
          <CandlestickSeries {...candlesAppearance} />
          <OHLCTooltip
            origin={[270, -15]}
            labelFill="#FFFFFF"
            textFill={(d) => handleColour(d)}
            fontSize={12.5}
            ohlcFormat={pricesDisplayFormat}
            changeFormat={changeDisplayFormat}
          />
          <EdgeIndicator
            itemType="last"
            orient="right"
            edgeAt="right"
            yAccessor={(d) => d.close}
            lineStroke="#FFFFFF"
            displayFormat={pricesDisplayFormat}
            fill={(d) => (d.close < d.open ? "#F70D1A" : "#52D017")}
            fontSize={13.5}
            fitToText={true}
          />
          <MouseCoordinateX
            at="bottom"
            orient="bottom"
            displayFormat={timeFormat("%Y-%m-%d %H:%M:%S")}
            fill="#303030"
          />
          <MouseCoordinateY
            displayFormat={pricesDisplayFormat}
            fill="#303030"
          />
          {false ? <ZoomButtons /> : null}
        </Chart>
        {selectedIndicators.includes("MACD") ? (
          <Chart
            id={2}
            origin={[0, chartHeight + spaceHeight]}
            yExtents={macdCalculator.accessor()}
            height={macdHeight}
          >
            <Macd
              macdCalculator={macdCalculator}
              macdAppearance={macdAppearance}
              isBottom={false}
            />
          </Chart>
        ) : null}
        {selectedIndicators.includes("Relative Strength Index") ? (
          <Chart
            id={3}
            origin={[0, chartHeight + macdHeight + 2 * spaceHeight]}
            yExtents={rsiCalculator.accessor()}
            height={rsiHeight}
          >
            <Rsi rsiCalculator={rsiCalculator} isBottom={false} />
          </Chart>
        ) : null}
        {selectedIndicators.includes("Stochastic Oscillator") ? (
          <Chart
            id={4}
            origin={[0, chartHeight + macdHeight + rsiHeight + 3 * spaceHeight]}
            yExtents={stoCalculator.accessor()}
            height={stoHeight}
          >
            <Stochastic stoCalculator={stoCalculator} isBottom={true} />
          </Chart>
        ) : null}
        <CrossHairCursor strokeDasharray="DashDot" strokeStyle="#4D4DFF" />
      </ChartCanvas>
    </React.Fragment>
  );
}

Quote.defaultProps = {
  type: "hybrid",
};

const currencyPairAlias = {
  EUR: "Euro",
  USD: "U.S. Dollar",
  GBP: "British Pound",
  CAD: "Canadian Dollar",
  JPY: "Japanese Yen",
  AUD: "Australian Dollar",
  NZD: "New Zealand Dollar",
  CHF: "Swiss Franc",
};

export default withSize({ style: { minHeight: 700 } })(
  withDeviceRatio()(Quote)
);
