import { Link, useLocation } from "react-router-dom";
import { useQueryParam, StringParam } from "use-query-params";

// names used for the search params
export const TELEMETRY_TYPE_PARAM_NAME = "tele";
export const CONFIGS_PARAM_NAME = "configs";
export const DESTINATIONS_PARAM_NAME = "dests";
export const PERIOD_PARAM_NAME = "period";

export function useQueryParamWrapper<Type>(
  searchParamName: string,
  defaultValue: Type
): readonly [
  searchParamsState: Type,
  setSearchParamsState: (newState: Type) => void
] {
  const [urlState, setUrlState] = useQueryParam(searchParamName, StringParam);

  const setState = (newState: Type) => {
    setUrlState(JSON.stringify(newState), "replaceIn");
  };
  const state = urlState ? (JSON.parse(urlState) as Type) : defaultValue;

  return [state, setState];
}

export interface SearchLinkProps {
  path: string;
  displayName: string;
}

export const SearchLink: React.FC<SearchLinkProps> = ({
  path,
  displayName,
}) => {
  const location = useLocation();
  return (
    <Link to={{ pathname: path, search: location.search }}>{displayName}</Link>
  );
};
