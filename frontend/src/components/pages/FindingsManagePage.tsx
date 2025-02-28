import {useState, useEffect} from "react"

import { activeProjectIdStore } from "@/lib/projectsStore";
import { useStore } from "@nanostores/react";

import { getProjectFindings } from "@/lib/findings";

import SimpleTable from "@/components/simple-table";
import PieChart from "@/components/pie-chart";

import CreateFindingDialog from "@/components/findings/create-finding-dialog";

import { type Payment, columns } from "../findings/data-table/columns";
import { DataTable } from "../findings/data-table/data-table";

const data = [
    {
        id: "728ed52f",
        amount: 100,
        status: "pending",
        email: "m@example.com",
    },
    {
        id: "728ed522",
        amount: 400,
        status: "failed",
        email: "b@example.com",
    },
    {
        id: "728ed525",
        amount: 300,
        status: "completed",
        email: "c@example.com",
    },
    {
        id: "728ed528",
        amount: 200,
        status: "pending",
        email: "d@example.com",
    },
    {
        id: "728ed52a",
        amount: 10,
        status: "pending",
        email: "e@example.com",
    },
    {
        id: "728ed52b",
        amount: 101,
        status: "pending",
        email: "f@example.com",
    },
]

export default function FindingsOverviewPage() {
    const $activeProjectId = useStore(activeProjectIdStore);
    const [findings, setFindings] = useState([]);
    const [scanConfigs, setScanConfigs] = useState({});

    useEffect(() => {
        if($activeProjectId) {
            getProjectFindings($activeProjectId).then(result => {
                if ("error" in result) {
                    return
                }
                console.log(result);
                setFindings(result);
            })
        }
    }, [$activeProjectId])

    return (
        <div className="mt-8 container mx-auto ">
            <h1>Findings</h1>
            <CreateFindingDialog />
            <div className="py-10">
                <DataTable columns={columns} data={findings} />
            </div>
        </div>
    );
}