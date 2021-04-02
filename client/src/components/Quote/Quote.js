import React from "react";
import { format } from "d3-format";
import { timeFormat } from "d3-time-format";
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
  Label,
} from "react-financial-charts";
import Rsi from "../Technical_Indicators/Rsi";
import Macd from "../Technical_Indicators/Macd";
import Stochastic from "../Technical_Indicators/Stochastic";
import BollingerBand from "../Technical_Indicators/Bollinger_Band";
import Ema from "../Technical_Indicators/Ema";

const candlesAppearance = {
  wickStroke: function fill(d) {
    return d.close > d.open ? "rgba(255, 7, 58, 1)" : "rgba(57, 255, 20, 1)";
  },
  fill: "#000000",
  stroke: function outline(d) {
    return d.close > d.open ? "rgba(255, 7, 58, 1)" : "rgba(57, 255, 20, 1)";
  },
  candleStrokeWidth: 1,
  width: 3,
};

function Quote(props) {
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

  const { type, height, width, ratio, data: initialData } = props;
  const calculatedData = macdCalculator(
    stoCalculator(rsiCalculator(bollingerCalculator(ema26(ema12(initialData)))))
  );
  const xScaleProvider = discontinuousTimeScaleProviderBuilder().inputDateAccessor(
    (d) => d.timestamp
  );
  const { data, xScale, xAccessor, displayXAccessor } = xScaleProvider(
    calculatedData
  );
  const max = xAccessor(data[data.length - 1]);
  const min = xAccessor(data.length > 200 ? data[data.length - 200] : data[0]);
  const xExtents = [min, max + 1];
  const pricesDisplayFormat = format(".5f");
  const margin = { left: 20, right: 70, top: 60, bottom: 30 };

  return (
    <React.Fragment>
      <ChartCanvas
        height={height}
        ratio={ratio}
        width={width}
        margin={margin}
        type={type}
        data={data}
        seriesName="MSFT"
        xAccessor={xAccessor}
        xScale={xScale}
        xExtents={xExtents}
        displayXAccessor={displayXAccessor}
        zoomAnchor={mouseBasedZoomAnchor}
      >
        {/*<Label
            x={(width - margin.left - margin.right) / 2}
            y={-5}
            fontSize="25"
            text={props.pair}
          />*/}
        <Chart
          id={1}
          yExtents={
            ((d) => [d.high, d.low],
            ema12.accessor(),
            ema26.accessor(),
            bollingerCalculator.accessor())
          }
          height={height / 3}
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
            strokeStyle="#FFFFFF"
            axisAt="right"
            orient="right"
            ticks={10}
            tickFormat={pricesDisplayFormat}
            tickLabelFill="#FFFFFF"
            gridLinesStrokeDasharray="Solid"
          />
          <BollingerBand bollingerCalculator={bollingerCalculator} />
          <Ema ema={[ema12, ema26]} />
          <CandlestickSeries {...candlesAppearance} />
          <EdgeIndicator
            itemType="last"
            orient="right"
            edgeAt="right"
            yAccessor={(d) => d.close}
            lineStroke="#FFFFFF"
            displayFormat={pricesDisplayFormat}
            fill={(d) => (d.close > d.open ? "#F70D1A" : "#52D017")}
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
        </Chart>
        <Chart
          id={2}
          origin={[0, height / 3 + 15]}
          yExtents={macdCalculator.accessor()}
          height={(height * 5) / 36}
        >
          <Macd
            macdCalculator={macdCalculator}
            macdAppearance={macdAppearance}
            isBottom={false}
          />
        </Chart>
        <Chart
          id={3}
          origin={[0, (height * 17) / 36 + 35]}
          yExtents={[0, 100]}
          height={(height * 5) / 36}
        >
          <Rsi rsiCalculator={rsiCalculator} isBottom={false} />
        </Chart>
        <Chart
          id={4}
          origin={[0, (height * 22) / 36 + 65]}
          yExtents={[0, 100]}
          height={(height * 5) / 36}
        >
          <Stochastic stoCalculator={stoCalculator} isBottom={true} />
        </Chart>
        <CrossHairCursor strokeDasharray="DashDot" strokeStyle="#4D4DFF" />
      </ChartCanvas>
    </React.Fragment>
  );
}

Quote.defaultProps = {
  type: "hybrid",
};

export default withSize({ style: { minHeight: 1000 } })(
  withDeviceRatio()(Quote)
);
