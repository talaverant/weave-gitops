import { CircularProgress } from "@material-ui/core";
import * as React from "react";
import styled from "styled-components";
import { useListFluxEvents } from "../hooks/events";
import { Event, ObjectReference } from "../lib/api/core/types.pb";
import { NoNamespace } from "../lib/types";
import Alert from "./Alert";
import DataTable from "./DataTable";
import Flex from "./Flex";
import Spacer from "./Spacer";
import Text from "./Text";
import Timestamp from "./Timestamp";

type Props = {
  className?: string;
  namespace?: string;
  involvedObject: ObjectReference;
};

function EventsTable({
  className,
  namespace = NoNamespace,
  involvedObject,
}: Props) {
  const { data, isLoading, error } = useListFluxEvents(
    namespace,
    involvedObject
  );

  if (isLoading) {
    return (
      <Flex wide center align>
        <CircularProgress />
      </Flex>
    );
  }

  if (error) {
    return (
      <Spacer padding="small">
        <Alert title="Error" message={error.message} severity="error" />
      </Spacer>
    );
  }

  return (
    <DataTable
      className={className}
      fields={[
        {
          value: (e: Event) => <Text capitalize>{e.reason}</Text>,
          label: "Reason",
        },
        { value: "message", label: "Message" },
        { value: "component", label: "Component" },
        {
          label: "Timestamp",
          value: (e: Event) => <Timestamp time={e.timestamp} />,
        },
      ]}
      rows={data.events}
    />
  );
}

export default styled(EventsTable).attrs({ className: EventsTable.name })`
  td {
    max-width: 1024px;
  }
`;
