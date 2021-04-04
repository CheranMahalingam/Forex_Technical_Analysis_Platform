import React from "react";
import CloseIcon from "@material-ui/icons/Close";
import DialogTitle from "@material-ui/core/DialogTitle";
import Chip from "@material-ui/core/Chip";
import TextField from "@material-ui/core/TextField";
import Autocomplete from "@material-ui/lab/Autocomplete";
import useWindowDimensions from "../../utils/useWindowDimensions";
import { Dialog, DialogContent } from "@material-ui/core";

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
            top: -windowHeight / 15,
          }}
        >
          <div>
            <Autocomplete
              multiple
              options={indicators}
              value={props.selectedIndicators}
              onChange={props.handleIndicatorsChange}
              filterSelectedOptions
              renderTags={(tagValue, getTagProps) =>
                tagValue.map((option, index) => (
                  <Chip label={option} {...getTagProps({ index })} />
                ))
              }
              renderInput={(params) => (
                <TextField
                  {...params}
                  variant="outlined"
                  label="Search"
                  margin="normal"
                />
              )}
              style={{
                position: "relative",
                margin: windowWidth / 60,
              }}
            />
          </div>
        </DialogContent>
      </Dialog>
    </React.Fragment>
  );
}

const indicators = [
  "Exponential Moving Average",
  "Relative Strength Index",
  "Stochastic Oscillator",
  "Bollinger Bands",
  "MACD",
];

export default SelectIndicatorModal;
