import _ from "lodash";
import * as React from "react";
import styled from "styled-components";
import { Condition } from "../lib/api/core/types.pb";
import Flex from "./Flex";
import Icon, { IconType } from "./Icon";
import Text from "./Text";

type Props = {
  className?: string;
  obj: Statusable;
  short?: boolean;
  suspended?: boolean;
};

interface Statusable {
  suspended?: boolean;
  conditions?: Condition[];
}

export function computeReady(conditions: Condition[]): boolean {
  const ready =
    _.find(conditions, { type: "Ready" }) ||
    // Deployment conditions work slightly differently;
    // they show "Available" instead of 'Ready'
    _.find(conditions, { type: "Available" });

  return ready?.status == "True";
}

export function computeMessage(obj: Statusable) {
  const readyCondition =
    _.find(obj.conditions, (c) => c.type === "Ready") ||
    _.find(obj.conditions, (c) => c.type === "Available");

  return readyCondition ? readyCondition.message : "unknown error";
}

function KubeStatusIndicator({ className, obj, short, suspended }: Props) {
  let readyText;
  let icon;
  if (suspended) {
    icon = IconType.SuspendedIcon;
    readyText = "Suspended";
  } else {
    const ready = computeReady(obj.conditions);
    readyText = ready ? "Ready" : "Not Ready";
    icon = readyText === "Ready" ? IconType.SuccessIcon : IconType.FailedIcon;
  }

  let text = computeMessage(obj);
  if (short || suspended) text = readyText;

  return (
    <Flex start className={className} align>
      <Icon size="base" type={icon} text={text} />
    </Flex>
  );
}

export default styled(KubeStatusIndicator).attrs({
  className: KubeStatusIndicator.name,
})`
  ${Icon} ${Text} {
    color: ${(props) => props.theme.colors.black};
    font-weight: 400;
  }
`;
