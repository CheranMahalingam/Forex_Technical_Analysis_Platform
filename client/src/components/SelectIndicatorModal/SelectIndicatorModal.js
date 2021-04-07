import React from "react";
import CloseIcon from "@material-ui/icons/Close";
import DialogTitle from "@material-ui/core/DialogTitle";
import Switch from "@material-ui/core/Switch";
import useWindowDimensions from "../../utils/useWindowDimensions";
import { Dialog, DialogContent, Typography } from "@material-ui/core";

function SelectIndicatorModal(props) {
  const { windowHeight, windowWidth } = useWindowDimensions();

  return (
    <React.Fragment>
      <Dialog open={props.openIndicators} onClose={props.handleIndicatorsClose}>
        <CloseIcon
          onClick={props.handleIndicatorsClose}
          style={{
            cursor: "pointer",
            marginTop: "4%",
            left: windowWidth / 3 - 10,
            position: "relative",
          }}
        />
        <DialogTitle>Indicators</DialogTitle>
        <DialogContent
          style={{
            position: "relative",
            width: windowWidth / 3,
            height: windowHeight / 3,
          }}
        >
          {props.indicators &&
            props.indicators.map((indicator) => {
              return (
                <div style={{ display: "flex" }}>
                  <Switch
                    name={indicator}
                    checked={props.selectedIndicators.includes(indicator)}
                    onChange={(e) =>
                      props.handleIndicatorsChange(e, props.pair)
                    }
                    color="primary"
                  />
                  <Typography variant="h6" style={{ marginLeft: "15%" }}>
                    {indicator}
                  </Typography>
                </div>
              );
            })}
        </DialogContent>
      </Dialog>
    </React.Fragment>
  );
}

export default SelectIndicatorModal;
