"use client"

import {useParams, useRouter} from "next/navigation";
import {useQuery} from "@tanstack/react-query";
import axios from "axios";
import {ProgressSpinner} from "primereact/progressspinner";
import {Card} from "primereact/card";
import React, {useState} from "react";
import Image from "next/image"
import {Button} from "primereact/button";

interface Candidate {
    ID: number;
    CreatedAt: string;
    UpdatedAt: string;
    DeletedAt: string | null;
    firstName: string;
    lastName: string;
    age: number;
    party: string;
    photo: string;
}

interface IGetCandidatesResponse {
    candidates: Candidate[];
    message: string;
}

export default function CandidatesPage() {
    const [isVoting, setIsVoting] = useState(false);
    const [errorVoting, setErrorVoting] = useState('');

    const params = useParams();
    const electionId = params.electionId;
    const router = useRouter();

    const { data, isLoading, error } = useQuery({
        queryKey: ['candidates', electionId],
        queryFn: async () => {
            const response = await axios.get<IGetCandidatesResponse>(
                `http://localhost:3000/elections/candidates/${electionId}`
            );
            return response.data;
        },
    });

    const handleVote = async (candidate: Candidate) => {
        const voterIDNP = localStorage.getItem('voterIDNP');

        if (!voterIDNP) {
            setErrorVoting('Voter IDNP not found. Please register first.');
            return;
        }

        setIsVoting(true);
        setErrorVoting('');

        try {
            const url = `http://localhost:3000/invoke?channelid=mychannel&chaincodeid=basic&function=CastVote&args=${voterIDNP}&args=${candidate.age}${candidate.party}`

            await axios.post(url);
        } catch (error) {
            setErrorVoting('Failed to cast vote. Please try again later.');
            console.error('Voting error:',error);
        } finally {
            setIsVoting(false);
        }
    }

    if (isLoading) {
        return (
            <div className="flex justify-center items-center min-h-screen">
                <ProgressSpinner />
            </div>
        )
    }

    if (error) {
        return (
            <div className="text-center text-red-500 p-4">
                Error loading candidates. Please try again later.
            </div>
        )
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
                    <h1 className="text-3xl font-bold">Election Candidates</h1>
                </div>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                {data?.candidates.map((candidate) => (
                    <Card
                        key={candidate.ID}
                        className="shadow-lg transform transistion-transform duration-200 hover:scale-105"
                        title={
                            <div className="text-xl font-bold text-center">
                                {candidate.firstName} {candidate.lastName}
                            </div>
                        }
                    >
                        <div className="space-y-4">
                            <div className="flex justify-center">
                                {
                                    candidate.photo ? (
                                        <div className="relative w-48 h-48 rounded-full overflow-hidden">
                                            <Image
                                                src={`data:image/jpeg;base64,${candidate.photo}`}
                                                alt={`${candidate.firstName} ${candidate.lastName}`}
                                                width={192}
                                                height={192}
                                                style={{
                                                    objectFit: 'cover'
                                                }}
                                                className="object-zoom"
                                            />
                                        </div>
                                    ) : (
                                        <div
                                            className="w-48 h-48 rounded-full bg-gray-200 flex items-center justify-center">
                                            <span className="text-4xl text-gray-400">
                                                {candidate.firstName[0]}{candidate.lastName[0]}
                                            </span>
                                        </div>
                                    )}
                            </div>
                            <div className="grid grid-cols-2 gap-4 text-sm">
                                <div>
                                    <p className="font-semibold">Age</p>
                                    <p>{candidate.age} years</p>
                                </div>
                                <div>
                                    <p className="font-semibold">Party</p>
                                    <p className={`${
                                        candidate.party === 'Democrat'
                                            ? 'text-blue-600'
                                            : candidate.party === 'Republican'
                                                ? 'text-red-600'
                                                : 'text-gray-600'
                                    }`}>
                                        {candidate.party}
                                    </p>
                                </div>
                            </div>
                            { errorVoting && (
                                <div className="text-red-500 text-sm text-center">
                                    {error}
                                </div>
                            )}
                            <div className="flex justify-center mt-4">
                                <Button
                                    icon="pi pi-check-circle"
                                    onClick={()=> handleVote(candidate)}
                                    disabled={isVoting}
                                    className={`w-full justify-content-center  ${isVoting ? 'opacity-50' : ''}`}
                                    label={isVoting ? 'Voting...' : 'Vote'}
                                />
                            </div>
                        </div>
                    </Card>
                ))}
            </div>
        </div>
    )
}