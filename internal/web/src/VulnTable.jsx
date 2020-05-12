import React, { useState, useEffect } from "react";

import Table from "@material-ui/core/Table";
import TableBody from "@material-ui/core/TableBody";
import TableCell from "@material-ui/core/TableCell";
import TableFooter from "@material-ui/core/TableFooter";
import TableHead from "@material-ui/core/TableHead";
import TablePagination from "@material-ui/core/TablePagination";
import TableRow from "@material-ui/core/TableRow";

import orange from "@material-ui/core/colors/orange";
import red from "@material-ui/core/colors/red";
import yellow from "@material-ui/core/colors/yellow";

export default function VulnTable({ target, vulnerabilities }) {
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const emptyRows =
    rowsPerPage -
    Math.min(rowsPerPage, vulnerabilities.length - page * rowsPerPage);

  const handleChangePage = (event, newPage) => {
    setPage(newPage);
  };

  const handleChangeRowsPerPage = event => {
    setRowsPerPage(parseInt(event.target.value, 10));
    setPage(0);
  };

  useEffect(() => {
    setRowsPerPage(Math.min(rowsPerPage, vulnerabilities.length));
  }, [vulnerabilities, rowsPerPage]);

  return (
    <Table stickyHeader aria-label="sticky table">
      <TableHead>
        <TableRow>
          <TableCell>Package Name</TableCell>
          <TableCell>Vulnerability ID</TableCell>
          <TableCell>Title</TableCell>
          <TableCell>Installed Version</TableCell>
          <TableCell>Fixed Version</TableCell>
          <TableCell>Severity</TableCell>
        </TableRow>
      </TableHead>
      <TableBody>
        {(rowsPerPage > 0
          ? vulnerabilities.slice(
              page * rowsPerPage,
              page * rowsPerPage + rowsPerPage
            )
          : vulnerabilities
        ).map((vuln, id) => (
          <TableRow key={id}>
            <TableCell>{vuln.PkgName}</TableCell>
            <TableCell style={{ maxWidth: "15em" }}>
              {vuln.VulnerabilityID}
              <sup style={{ display: "inline" }}>
                {(vuln.References || []).map((ref, i) => (
                  <>
                    <a href={ref} style={{ color: "white" }}>
                      {i}
                    </a>{" "}
                  </>
                ))}
              </sup>
            </TableCell>
            <TableCell>{vuln.Title || "N/A"}</TableCell>
            <TableCell>{vuln.InstalledVersion}</TableCell>
            <TableCell>{vuln.FixedVersion || "N/A"}</TableCell>
            <TableCell
              style={{
                color:
                  vuln.Severity === "HIGH" || vuln.Severity === "CRITICAL"
                    ? red[200]
                    : vuln.Severity === "MEDIUM"
                    ? orange[200]
                    : yellow[200]
              }}
            >
              {vuln.Severity}
            </TableCell>
          </TableRow>
        ))}
        {emptyRows > 0 && (
          <TableRow style={{ height: 53 * emptyRows }}>
            <TableCell colSpan={6} />
          </TableRow>
        )}
      </TableBody>
      <TableFooter>
        <TableRow>
          <TablePagination
            rowsPerPageOptions={[5, 10, 25, { label: "All", value: -1 }]}
            colSpan={6}
            count={vulnerabilities.length}
            rowsPerPage={rowsPerPage}
            page={page}
            SelectProps={{
              inputProps: { "aria-label": "rows per page" },
              native: true
            }}
            onChangePage={handleChangePage}
            onChangeRowsPerPage={handleChangeRowsPerPage}
          />
        </TableRow>
      </TableFooter>
    </Table>
  );
}
