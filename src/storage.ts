// handles GSC and local file interactions

import { Storage } from "@google-cloud/storage";
import fs from 'fs';
import ffmpeg from 'fluent-ffmpeg';

const storage = new Storage();

const rawVideoBucketName = "rk_sample_vid.mp4";
const processsedVideoBucketName = "rk_processed_sample_vid.mp4";

const localRawVideoPath =  "./raw-videos";
const localprocesssedVideoPath = "./processed-videos";

//creates local dir for processsed and raw
export function setupDirectories(){

}


/**
 * @param rawVideoName - The name of the file to convert from {@link localRawVideoPath}.
 * @param processedVideoName - The name of the file to convert to {@link localProcessedVideoPath}.
 * @returns A promise that resolves when the video has been converted.
 */
export function convertVideo(rawVideoName: string, processsedVideoName: string){
    return new Promise<void>((resolve, reject) => {
        ffmpeg(`${localRawVideoPath}/${rawVideoName}`)
            .outputOption('-vf', 'scale=-1:360') //360p
            .on('end', () => {
                console.log("Processing finished succesfully.");
                resolve();
            })
            .on('error', (err) => {
                console.log(`An error occured: ${err.message}`);
                reject(err);
            })
            .save(`${localprocesssedVideoPath}/${processsedVideoName}`)
            })
    
}
/**
 * @param fileName - The name of the file to download from the 
 * {@link rawVideoBucketName} bucket into the {@link localRawVideoPath} folder.
 * @returns A promise that resolves when the file has been downloaded.
 */

export async function downloadRawVideo(fileName: string){
    await storage.bucket(rawVideoBucketName)
        .file(fileName)
        .download({ destination: `${localRawVideoPath}/${fileName}` });
    
    console.log(
        `gs.//${rawVideoBucketName}/${fileName} downloaded to ${localRawVideoPath}/${fileName}.` 
    )
}

/**
 * @param fileName - The name of the file to upload from the 
 * {@link localProcessedVideoPath} folder into the {@link processedVideoBucketName}.
 * @returns A promise that resolves when the file has been uploaded.
 */
export async function uploadProcessedVideo(fileName:string){
    const bucket = storage.bucket(processsedVideoBucketName);

    await bucket.upload(`${localprocesssedVideoPath}/${fileName}`, {
        destination: fileName
    });
    console.log(`gs://${rawVideoBucketName}/${fileName} downloaded to ${localRawVideoPath}/${fileName}.`)
    await bucket.file(fileName).makePublic();

}