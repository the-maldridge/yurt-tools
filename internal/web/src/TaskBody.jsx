import React, { useState } from "react";

import red from "@material-ui/core/colors/red";
import green from "@material-ui/core/colors/green";
import LockIcon from "@material-ui/icons/Lock";
import UpdateIcon from "@material-ui/icons/Update";
import Divider from "@material-ui/core/Divider";
import Typography from "@material-ui/core/Typography";
import Card from "@material-ui/core/Card";
import CardContent from "@material-ui/core/CardContent";
import CardActions from "@material-ui/core/CardActions";
import Button from "@material-ui/core/Button";
import Grid from "@material-ui/core/Grid";
import Skeleton from "@material-ui/lab/Skeleton";

import VulnTable from "./VulnTable";
import VulnSummary from "./VulnSummary";

const FakeCard = () => (
  <Skeleton height={185} variant="rect" animation="wave" />
);

export default function TaskBody({ currentTask, job, group, task }) {
  const [vulnSummaryOpen, setVulnSummaryOpen] = useState(false);

  const vulnerable =
    currentTask && currentTask.trivy
      ? currentTask.trivy
          .map(target => !!target.Vulnerabilities)
          .reduce((acc, cur) => cur || acc, false)
      : null;

  return (
    <Grid container spacing={3}>
      <Grid item xs={12}>
        <Typography variant="h5" color="inherit">
          {`${job} » ${group} » ${task}`}
        </Typography>
      </Grid>
      <Divider />
      <Grid item xs={4}>
        {currentTask ? (
          <Card>
            <CardContent>
              <Typography color="textSecondary" gutterBottom>
                <>docker</>
              </Typography>
              <Typography variant="h5" component="h2">
                {currentTask.metadata.Docker.Image}
              </Typography>
              <Typography color="textSecondary">
                {currentTask.metadata.Docker.Tag}
              </Typography>
              <Typography variant="body2" component="p">
                {currentTask.metadata.Docker.Owner || "Docker Official Image"}
              </Typography>
            </CardContent>
            <CardActions>
              <Button size="small" href={currentTask.metadata.URL}>
                Registry
              </Button>
            </CardActions>
          </Card>
        ) : (
          <FakeCard />
        )}
      </Grid>
      <Grid item xs={4}>
        {currentTask ? (
          <Card>
            <CardContent>
              <Typography color="textSecondary" gutterBottom>
                <>security</>
              </Typography>
              <Typography variant="h5" component="h2" style={{ color: green }}>
                {vulnerable !== null ? (
                  vulnerable ? (
                    <span>
                      Vulnerabilities found{" "}
                      <LockIcon
                        fontSize="small"
                        style={{ color: red[200], verticalAlign: "middle" }}
                      />
                    </span>
                  ) : (
                    <span>
                      No known vulnerabilities{" "}
                      <LockIcon
                        fontSize="small"
                        style={{ color: green[200], verticalAlign: "middle" }}
                      />
                    </span>
                  )
                ) : (
                  "No vulnerability data"
                )}
              </Typography>
              <Typography color="textSecondary">
                <span>
                  {vulnerable !== null
                    ? Object.keys(currentTask.trivy).length
                    : 0}{" "}
                  targets scanned
                </span>
              </Typography>
              <Typography variant="body2" component="p">
                <span>
                  {vulnerable !== null
                    ? Object.values(currentTask.trivy)
                        .map(target =>
                          target.Vulnerabilities
                            ? Object.keys(target.Vulnerabilities).length
                            : 0
                        )
                        .reduce((acc, cur) => acc + cur, 0)
                    : 0}{" "}
                  vulnerabilities discovered
                </span>
              </Typography>
            </CardContent>
            <CardActions>
              <Button size="small" onClick={() => setVulnSummaryOpen(true)}>
                View Summary
              </Button>
            </CardActions>
            <VulnSummary
              vulnerabilities={currentTask.trivy}
              open={vulnSummaryOpen}
              onClose={() => setVulnSummaryOpen(false)}
            />
          </Card>
        ) : (
          <FakeCard />
        )}
      </Grid>
      <Grid item xs={4}>
        {currentTask ? (
          <Card>
            <CardContent>
              <Typography color="textSecondary" gutterBottom>
                version
              </Typography>
              <Typography variant="h5" component="h2" style={{ color: green }}>
                {currentTask.versions && !currentTask.versions.NonComparable ? (
                  currentTask.versions.UpToDate ? (
                    <span>
                      Up to date{" "}
                      <UpdateIcon
                        fontSize="small"
                        style={{ color: green[200], verticalAlign: "middle" }}
                      />
                    </span>
                  ) : (
                    <span>
                      Out of date{" "}
                      <UpdateIcon
                        fontSize="small"
                        style={{ color: red[200], verticalAlign: "middle" }}
                      />
                    </span>
                  )
                ) : (
                  "Unknown version"
                )}
              </Typography>
              <Typography color="textSecondary">
                {currentTask.versions && !currentTask.versions.NonComparable
                  ? currentTask.versions.UpToDate
                    ? currentTask.versions.Current
                    : `${currentTask.versions.Available.length} newer versions available`
                  : "No Data"}
              </Typography>
              <Typography variant="body2" component="p">
                {currentTask.versions && currentTask.versions.Available
                  ? `Latest: ${currentTask.versions.Available[0]}`
                  : currentTask.versions
                  ? "Latest version"
                  : "No Data"}
              </Typography>
            </CardContent>
            <CardActions>
              <Button
                size="small"
                href={`${currentTask.metadata.URL}?tab=tags`}
              >
                Tags
              </Button>
            </CardActions>
          </Card>
        ) : (
          <FakeCard />
        )}
      </Grid>
      {currentTask && currentTask.trivy ? (
        currentTask.trivy
          .filter(target => target.Vulnerabilities)
          .map(({ Target, Vulnerabilities }) => (
            <React.Fragment key={Target}>
              <Typography variant="h6" style={{ paddingLeft: "15px" }}>
                {Target}
              </Typography>
              <VulnTable target={Target} vulnerabilities={Vulnerabilities} />
            </React.Fragment>
          ))
      ) : (
        <></>
      )}
    </Grid>
  );
}
