"use client"

import {useParams, useRouter} from "next/navigation";
import React, {useEffect, useState} from "react";
import axios from "axios";
import "primereact/resources/themes/lara-light-cyan/theme.css";
import {DataTable} from "primereact/datatable";
import {Column} from "primereact/column";
import {Button} from "primereact/button";
import {Badge} from "primereact/badge";

export default function ResultsPage()  {
    const router = useRouter();
    const params = useParams();
    const electionId = params.electionId;

    const [results, setResults] = useState<any[]>([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        if (electionId) {
            const fetchResults = async () => {
                try {
                    const response = await axios.get(`http://localhost:3000/query`, {
                        params: {
                            channelid: 'mychannel',
                            chaincodeid: 'basic',
                            function: 'GetCandidatesByElection',
                            args: electionId,
                        },
                    });
                    const parsedData =
                        typeof response.data === 'string'
                            ? JSON.parse(response.data.replace('Response: ', ''))
                            : response.data;
                    setResults(parsedData)
                } catch (error) {
                    console.error('Error fetching results:', error);
                } finally {
                    setLoading(false);
                }
            };

            fetchResults();
        }
    }, [electionId]);

    if (loading) {
        return <p>Loading...</p>
    }

    if (results.length === 0) {
        return <p>No results found.</p>;
    }

    return (
        <div className="p-6">
            <div className="flex justify-content-start">
                <Button
                    style={{height: 40}}
                    className="mt-3"
                    severity="secondary"
                    icon="pi pi-angle-left"
                    label="Elections"
                    raised
                    onClick={() => router.push('/election')}
                />
            </div>
            <div className="flex justify-content-center">
                <h1 className="text-3xl font-bold">Election Results</h1>
            </div>
            <DataTable className="p-datatable p-component" value={results}>
                <Column field="id" header="ID"></Column>
                <Column field="name" header="Name"></Column>
                <Column field="electionID" header="Election ID"></Column>
                <Column field="votes" header="Votes"></Column>
                <Column
                    header="Winner"
                    body={(rowData) => {
                        const maxVotes = Math.max(...results.map((candidate) => candidate.votes));
                        return rowData.votes === maxVotes ? (
                            <Badge value="Winner" severity="success" />
                        ) : null;
                    }}
                ></Column>
            </DataTable>
        </div>
    );
}