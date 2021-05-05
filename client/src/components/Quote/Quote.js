import React, { useState } from "react";
import Button from "@material-ui/core/Button";
import AddIcon from "@material-ui/icons/Add";
import TimelineIcon from "@material-ui/icons/Timeline";
import RemoveIcon from "@material-ui/icons/Remove";
import { Chart, ChartCanvas } from "@react-financial-charts/core";
import { discontinuousTimeScaleProviderBuilder } from "@react-financial-charts/scales";
import { withDeviceRatio, withSize } from "@react-financial-charts/utils";
import {
  ema,
  bollingerBand,
  rsi,
  macd,
  stochasticOscillator,
} from "@react-financial-charts/indicators";
import {
  mouseBasedZoomAnchor,
  Label,
  CrossHairCursor,
} from "react-financial-charts";
import Rsi from "../Technical_Indicators/Rsi";
import Macd from "../Technical_Indicators/Macd";
import Stochastic from "../Technical_Indicators/Stochastic";
import SelectIndicatorModal from "../SelectIndicatorModal/SelectIndicatorModal";
import ExchangeRate from "../ExchangeRate/ExchangeRate";

function Quote(props) {
  const { type, width, ratio, data: initialData } = props;
  const height = 750;

  const [openIndicators, setOpenIndicators] = useState(false);

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
  const min = xAccessor(data.length > 100 ? data[data.length - 100] : data[0]);
  const xExtents = [min, max + 1];
  const margin = { left: 20, right: 80, top: 50, bottom: 30 };
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

  return (
    <React.Fragment>
      <div
        style={{
          position: "absolute",
          top: props.quoteHeight[props.pair][0],
        }}
      >
        <Button
          variant="outlined"
          style={{
            position: "relative",
            left: (5 * width) / 6,
            top: margin.top,
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
        <Button
          variant="outlined"
          style={{
            position: "relative",
            left: ((5 * width) / 6) - 300,
            top: margin.top,
            zIndex: 7,
            color: "#FFFFFF",
            borderColor: "#FFFFFF",
          }}
          onClick={() => props.handleInferenceSubscribe(props.pair)}
          startIcon={
            props.isSubscribedToInference ? <RemoveIcon /> : <TimelineIcon />
          }
          size="small"
        >
          View Forecast
        </Button>
        <SelectIndicatorModal
          openIndicators={openIndicators}
          handleIndicatorsClose={handleIndicatorsClose}
          selectedIndicators={props.selectedIndicators}
          handleIndicatorsChange={props.handleIndicatorsChange}
          indicators={indicators}
          pair={props.pair}
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
          <Label
            x={0}
            y={-15}
            fontSize={18}
            text={pairLabel}
            textAlign="left"
          />
          <Chart
            id={1}
            yExtents={
              ((d) => [d.high + 0.005, d.low - 0.005],
              (d) => [d.inference + 0.005, d.inference - 0.005],
              ema12.accessor(),
              ema26.accessor(),
              bollingerCalculator.accessor())
            }
            height={props.indicatorHeight[0]}
          >
            <ExchangeRate
              selectedIndicators={props.selectedIndicators}
              bollingerCalculator={bollingerCalculator}
              ema={[ema12, ema26]}
              handleColour={handleColour}
              isBottom={
                Math.max(...props.indicatorHeight) === props.indicatorHeight[0]
              }
              moveRight={(3 * width) / 13}
              isSubscribedToInference={props.isSubscribedToInference}
            />
          </Chart>
          {props.selectedIndicators.includes("MACD") ? (
            <Chart
              id={2}
              origin={[0, props.indicatorHeight[1] - props.macdHeight]}
              yExtents={macdCalculator.accessor()}
              height={props.macdHeight}
            >
              <Macd
                macdCalculator={macdCalculator}
                macdAppearance={macdAppearance}
                isBottom={
                  Math.max(...props.indicatorHeight) ===
                  props.indicatorHeight[1]
                }
              />
            </Chart>
          ) : null}
          {props.selectedIndicators.includes("Relative Strength Index") ? (
            <Chart
              id={3}
              origin={[0, props.indicatorHeight[2] - props.rsiHeight]}
              yExtents={rsiCalculator.accessor()}
              height={props.rsiHeight}
            >
              <Rsi
                rsiCalculator={rsiCalculator}
                isBottom={
                  Math.max(...props.indicatorHeight) ===
                  props.indicatorHeight[2]
                }
              />
            </Chart>
          ) : null}
          {props.selectedIndicators.includes("Stochastic Oscillator") ? (
            <Chart
              id={4}
              origin={[0, props.indicatorHeight[3] - props.stoHeight]}
              yExtents={stoCalculator.accessor()}
              height={props.stoHeight}
            >
              <Stochastic
                stoCalculator={stoCalculator}
                isBottom={
                  Math.max(...props.indicatorHeight) ===
                  props.indicatorHeight[3]
                }
              />
            </Chart>
          ) : null}
          <CrossHairCursor strokeDasharray="DashDot" strokeStyle="#4D4DFF" />
        </ChartCanvas>
      </div>
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

const indicators = [
  "Exponential Moving Average",
  "Relative Strength Index",
  "Stochastic Oscillator",
  "Bollinger Bands",
  "MACD",
];

export default withSize({ style: { minHeight: 800 } })(
  withDeviceRatio()(Quote)
);
