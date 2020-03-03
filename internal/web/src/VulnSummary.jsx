import React from "react";
import Container from "@material-ui/core/Container";
import ListItem from "@material-ui/core/ListItem";
import DialogTitle from "@material-ui/core/DialogTitle";
import Dialog from "@material-ui/core/Dialog";
import { red, orange, yellow } from "@material-ui/core/colors";

export default function VulnSummary({ onClose, vulnerabilities, open }) {
  return (
    <Dialog onClose={onClose} aria-labelledby="simple-dialog-title" open={open}>
      <DialogTitle id="simple-dialog-title">Vulnerability Summary</DialogTitle>
      <Container>
        {vulnerabilities ? (
          vulnerabilities.map(({ Target, Vulnerabilities }) => {

            const bySeverity = Vulnerabilities ? Vulnerabilities.reduce(
              (acc, cur) => ({
                ...acc,
                [cur.Severity]: (acc[cur.Severity] || 0) + 1
              }),
              {}
            ) : ({None: 0});
            return (
              <>
                <b>{Target}</b>
                {Object.entries(bySeverity).map(([severity, count]) => (
                  <div style={{paddingBottom: "5px"}}>
                    <span
                      style={{
                        color:
                          severity === "HIGH" || severity === "CRITICAL"
                            ? red[200]
                            : severity === "MEDIUM"
                            ? orange[200]
                            : yellow[200]
                      }}
                    >
                      {severity}
                    </span>
                    : {count}
                  </div>
                ))}
              </>
            );
          })
        ) : (
          <ListItem>No targets scanned</ListItem>
        )}
      </Container>
    </Dialog>
  );
}
