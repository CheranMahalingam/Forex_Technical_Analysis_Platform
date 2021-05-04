import React from "react";
import { LineSeries, LineSeriesProps } from "@react-financial-charts/series";

function Inference(props) {
  return (
    <React.Fragment>
      <LineSeries
        yAccessor={props.yAccessor}
        strokeStyle="#FFFFFF"
        connectNulls={true}
      />
    </React.Fragment>
  );
}

export default Inference;
