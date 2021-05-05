import React from "react";
import { LineSeries} from "@react-financial-charts/series";

function Inference(props) {
  return (
    <React.Fragment>
      <LineSeries
        yAccessor={props.yAccessor}
        strokeStyle="#FFFFFF"
      />
    </React.Fragment>
  );
}

export default Inference;
